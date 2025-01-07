package api

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (a App) PasswordLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	r.ParseForm()

	email, exists := r.Form["email"]
	password, exists2 := r.Form["password"]

	_, emailErr := mail.ParseAddress(email[0])
	if !exists || !exists2 || emailErr != nil || len(password[0]) < 8 {
		// Some field is empty
		http.Error(w, "Bad form fields", http.StatusBadRequest)
		return
	}

	// we return InvalidCredentials either way of not finding email or password mismatch
	// to prevent malicious actors from guessing correct emails
	user, err := a.db.GetUserByEmail(r.Context(), email[0])
	if errors.Is(pgx.ErrNoRows, err) {
		// No user found
		response := ApiResponse[db.User]{ErrCode: util.Of(int(InvalidCredentials))}
		JSONError(w, response, http.StatusNotFound)
		return
	}

	if err != nil {
		internalServerError(w)
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
		internalServerError(w)
		return
	}

	http.SetCookie(w, cookie)

	response := ApiResponse[db.User]{Data: []db.User{user}}
	JSON(w, response)
}

func (a App) PasswordSignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	r.ParseForm()

	username, exists := r.Form["username"]
	email, exists2 := r.Form["email"]
	password, exists3 := r.Form["password"]

	_, emailErr := mail.ParseAddress(email[0])
	if !exists || !exists2 || !exists3 || emailErr != nil || len(username[0]) < 3 || len(password[0]) < 8 {
		// Some field is empty
		http.Error(w, "Bad form fields", http.StatusBadRequest)
		return
	}

	_, err := a.db.GetUserByEmail(r.Context(), email[0])
	if err == nil {
		// a user with this email already exists
		response := ApiResponse[db.User]{ErrCode: util.Of(int(EmailExists))}
		JSONError(w, response, http.StatusOK)
		return
	} else if !errors.Is(pgx.ErrNoRows, err) {
		internalServerError(w)
		return
	}

	passwordHash, err := auth.GenerateHashString(password[0])
	if err != nil {
		internalServerError(w)
		return
	}

	user, err := a.db.CreateUser(r.Context(), db.CreateUserParams{Username: username[0], Email: email[0], PasswordHash: &passwordHash})
	if err != nil {
		internalServerError(w)
		return
	}

	cookie, err := a.createSessionCookie(r.Context(), user.ID)
	if err != nil {
		internalServerError(w)
		return
	}

	http.SetCookie(w, cookie)

	response := ApiResponse[db.CreateUserRow]{Data: []db.CreateUserRow{user}}
	JSON(w, response)
}
