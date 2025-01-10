package api

import (
	"log"
	"net/http"

	"github.com/spectre-xenon/lumina-chat/internal/db"
)

func (a App) LogoutSessionHandler(w http.ResponseWriter, r *http.Request, session db.Session) {
	dbErr := a.db.DeleteSession(r.Context(), session.SessionToken)
	if dbErr != nil {
		log.Printf("Database error: %s\n", dbErr)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	cookie := a.createEmptyCookie()

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a App) LogoutAllSessionsHandler(w http.ResponseWriter, r *http.Request, session db.Session) {
	dbErr := a.db.DeleteSessionsByUser(r.Context(), session.UserID)
	if dbErr != nil {
		log.Printf("Database error: %s\n", dbErr)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	cookie := a.createEmptyCookie()

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
