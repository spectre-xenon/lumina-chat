package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/hash"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type userData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newOAuthClient() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  strings.Join([]string{os.Getenv("ORIGIN"), "/v1/auth/callback/google"}, ""),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}

func redirectToLoginWithCode(w http.ResponseWriter, r *http.Request, code int) {
	url := strings.Join([]string{"/login?errcode=", strconv.Itoa(code)}, "")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (a *App) oauthLoginHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := a.ValidateSession(r)
	if ok {
		redirectToLoginWithCode(w, r, AlreadyLoggedIn)
		return
	}

	state := hash.RandString(10)

	client := newOAuthClient()
	url := client.AuthCodeURL(string(state))

	stateCookie := &http.Cookie{
		Name:     "state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 10,
	}

	http.SetCookie(w, stateCookie)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (a *App) oauthSignupHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state == "" || code == "" {
		http.Error(w, "Bad OAuth request", http.StatusBadRequest)
		return
	}

	prevState, err := r.Cookie("state")
	if errors.Is(http.ErrNoCookie, err) {
		redirectToLoginWithCode(w, r, NoStateCookie)
		return
	}

	if state != prevState.Value {
		redirectToLoginWithCode(w, r, MismatchedState)
		return
	}

	client := newOAuthClient()

	token, err := client.Exchange(r.Context(), code, oauth2.AccessTypeOffline)
	if err != nil {
		log.Printf("Error exchanging oauth code: %s\n", err)
		redirectToLoginWithCode(w, r, InternalServerError)
		return
	}

	url := strings.Join([]string{"https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=", token.AccessToken}, "")
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting user oauth info: %s\n", err)
		redirectToLoginWithCode(w, r, InternalServerError)
		return
	}

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading user oauth info from body: %s\n", err)
		redirectToLoginWithCode(w, r, InternalServerError)
		return
	}

	var userData userData
	marshalErr := json.Unmarshal(bodyData, &userData)
	if marshalErr != nil {
		log.Printf("Error unmarshaling user oauth info: %s\n", err)
		redirectToLoginWithCode(w, r, InternalServerError)
		return
	}

	var userID uuid.UUID
	// Create newUser if no newUser with that email was found
	oldUser, err := a.db.GetUserByEmail(r.Context(), userData.Email)
	if errors.Is(pgx.ErrNoRows, err) {
		newUser, err := a.db.CreateUser(r.Context(),
			db.CreateUserParams{
				Username:     userData.Name,
				Email:        userData.Email,
				PasswordHash: nil,
			})
		if err != nil {
			log.Printf("Database error: %s\n", err)
			redirectToLoginWithCode(w, r, InternalServerError)
			return
		}

		userID = newUser.ID
	} else {
		userID = oldUser.ID
	}

	cookie, err := a.createSessionCookie(r.Context(), userID)
	if err != nil {
		log.Printf("Error creating session cookie: %s\n", err)
		redirectToLoginWithCode(w, r, InternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
