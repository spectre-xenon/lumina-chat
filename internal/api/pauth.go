package api

import (
	"errors"
	"log"
	"net/http"
	"net/mail"

	"github.com/jackc/pgx/v5"
	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/hash"
	"github.com/spectre-xenon/lumina-chat/internal/util"
)

func (a *App) passwordLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	r.ParseForm()

	email, exists := r.Form["email"]
	password, exists2 := r.Form["password"]

	if !exists || !exists2 {
		// Some field is empty
		http.Error(w, "Bad form fields", http.StatusBadRequest)
		return
	}
	_, emailErr := mail.ParseAddress(email[0])

	if emailErr != nil || len(password[0]) < 8 {
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
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	if user.PasswordHash == nil {
		// NoPassword means the user signed up with an oauth provider
		response := ApiResponse[db.User]{ErrCode: util.Of(NoPassword)}
		JSONError(w, response, http.StatusOK)
		return
	}

	resp := a.workerPool.Submit(func() any {
		return hash.CompareHashStrings(password[0], *user.PasswordHash)
	})

	result := <-resp

	if !result.(bool) {
		// Invalid Credientials
		response := ApiResponse[db.User]{ErrCode: util.Of(int(InvalidCredentials))}
		JSONError(w, response, http.StatusOK)
		return
	}

	cookie, err := a.createSessionCookie(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error creating session cookie: %s\n", err)
		internalServerError(w)
		return
	}

	http.SetCookie(w, cookie)

	// Don't send the password hash back
	user.PasswordHash = nil
	response := ApiResponse[db.User]{Data: []db.User{user}}
	JSON(w, response)
}

func (a *App) passwordSignupHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	var hashErr error
	resp := a.workerPool.Submit(func() any {
		passwordHash, err := hash.GenerateHashString(password[0])
		hashErr = err
		return passwordHash
	})

	result := <-resp

	if hashErr != nil {
		log.Printf("Error generating password hash: %s\n", err)
		internalServerError(w)
		return
	}

	user, err := a.db.CreateUser(r.Context(), db.CreateUserParams{Username: username[0], Email: email[0], PasswordHash: result.(*string)})
	if err != nil {
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	cookie, err := a.createSessionCookie(r.Context(), user.ID)
	if err != nil {
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	http.SetCookie(w, cookie)

	response := ApiResponse[db.CreateUserRow]{Data: []db.CreateUserRow{user}}
	JSON(w, response)
}
