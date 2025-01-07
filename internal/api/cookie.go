package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spectre-xenon/lumina-chat/internal/db"
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
