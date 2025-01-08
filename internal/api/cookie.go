package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spectre-xenon/lumina-chat/internal/db"
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
