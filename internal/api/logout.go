package api

import (
	"net/http"

	"github.com/spectre-xenon/lumina-chat/internal/db"
)

func (a App) LogoutSessionHandler(w http.ResponseWriter, r *http.Request, session db.Session) {
	sessionToken := SafeParseSessionToken(r)

	dbErr := a.db.DeleteSession(r.Context(), sessionToken)
	if dbErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	cookie := a.createEmptyCookie()

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a App) LogoutAllSessionsHandler(w http.ResponseWriter, r *http.Request, session db.Session) {
	sessionToken := SafeParseSessionToken(r)

	session, err := a.db.GetSession(r.Context(), sessionToken)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	dbErr := a.db.DeleteSessionsByUser(r.Context(), session.UserID)
	if dbErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	cookie := a.createEmptyCookie()

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
