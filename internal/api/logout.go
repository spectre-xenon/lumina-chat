package api

import "net/http"

func (a App) LogoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken := SafeParseSessionToken(r)

	dbErr := a.db.DeleteSession(r.Context(), sessionToken)
	if dbErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	cookie := a.createEmptyCookie()

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a App) LogoutAllSessionsHandler(w http.ResponseWriter, r *http.Request) {
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
