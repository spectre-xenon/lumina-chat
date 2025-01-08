package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/util"
)

func SafeParseSessionToken(r *http.Request) (sessionToken uuid.UUID) {
	// Cookie is always available as we check before handing the request to here
	sessionCookie, _ := r.Cookie("session")
	sessionToken, _ = uuid.Parse(sessionCookie.Value)
	return
}

func (a App) createSessionCookie(ctx context.Context, userID uuid.UUID) (cookie *http.Cookie, err error) {
	// expires after 10 days
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

func (a App) createEmptyCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}
}

type ValidateSession func(r *http.Request) (session *db.Session, ok bool)

func (a App) ValidateSession(r *http.Request) (session *db.Session, ok bool) {
	sessionCookie, err := r.Cookie("session")
	if errors.Is(http.ErrNoCookie, err) {
		return nil, false
	}

	parsedToken, err := uuid.Parse(sessionCookie.Value)
	if errors.Is(http.ErrNoCookie, err) {
		return nil, false
	}

	dbSession, err := a.db.GetSession(r.Context(), parsedToken)
	if err != nil {
		return nil, false
	}

	now := time.Now()
	if now.After(dbSession.ExpiresAt) || now.Equal(dbSession.ExpiresAt) {
		return nil, false
	}

	return &dbSession, true
}

func (a App) WithAuth(next HanlderWithSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := a.ValidateSession(r)
		if !ok {
			http.Error(w, "UnAuthorized", http.StatusUnauthorized)
			return
		}

		next(w, r, *session)
	}
}

func (a App) WithNoAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := a.ValidateSession(r)
		if ok {
			// No user found
			response := ApiResponse[db.User]{ErrCode: util.Of(int(AlreadyLoggedIn))}
			JSON(w, response)
			return
		}

		next(w, r)
	}
}
