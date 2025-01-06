package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spectre-xenon/lumina-chat/internal/auth"
	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/util"
)

func (a App) createSessionCookie(ctx context.Context, userID uuid.UUID) (cookie *http.Cookie, err error) {
	expiresAt := time.Now().AddDate(0, 0, 10)
	sessionToken, err := a.db.CreateSession(ctx, db.CreateSessionParams{UserID: userID, ExpiresAt: expiresAt})

	cookie = &http.Cookie{
		Name:     "session",
		Value:    sessionToken.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  expiresAt,
	}

	return cookie, err
}

func (a App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	r.ParseForm()

	email, exists := r.Form["email"]
	password, exists2 := r.Form["password"]

	if !exists || !exists2 {
		// Some field is empty
		http.Error(w, "Empty form fields", http.StatusBadRequest)
		return
	}

	// we return InvalidCredentials either way of not finding email or password mismatch
	// to prevent malicious actors from guessing correct emails
	user, err := a.db.GetUserByEmail(r.Context(), email[0])
	if err != nil {
		// No user found
		response := ApiResponse[db.User]{ErrCode: util.Of(int(InvalidCredentials))}
		JSONError(w, response, http.StatusNotFound)
		return
	}

	if user.PasswordHash == nil {
		// NoPassword means the user signed up with an oauth provider
		response := ApiResponse[db.User]{ErrCode: util.Of(NoPassword)}
		JSONError(w, response, http.StatusOK)
		return
	}

	if !auth.CompareHashStrings(password[0], *user.PasswordHash) {
		// Invalid Credientials
		response := ApiResponse[db.User]{ErrCode: util.Of(int(InvalidCredentials))}
		JSONError(w, response, http.StatusOK)
		return
	}

	cookie, err := a.createSessionCookie(r.Context(), user.ID)
	if err != nil {
		response := ApiResponse[db.User]{ErrCode: nil}
		JSONError(w, response, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	response := ApiResponse[db.User]{Data: []db.User{user}}
	JSON(w, response)
}
