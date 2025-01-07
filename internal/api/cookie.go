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

func (a App) LogoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie is always available as we check before handing the request to here
	sessionCookie, _ := r.Cookie("session")

	uuid, err := uuid.Parse(sessionCookie.Value)
	if err != nil {
		http.Error(w, "Bad Session Cookie", http.StatusBadRequest)
		return
	}

	dbErr := a.db.DeleteSession(r.Context(), uuid)
	if dbErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	clearedCookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}

	http.SetCookie(w, clearedCookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
