package api

import (
	"net/http"

	"github.com/spectre-xenon/lumina-chat/internal/db"
)

type App struct {
	db  *db.Queries
	Mux *http.ServeMux
}

func New(db *db.Queries, mux *http.ServeMux) App {
	return App{db, mux}
}

func (a App) HandleFunc(pattern string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(pattern, handler)
}

func (a App) HandleFuncWithAuth(pattern string, handler HanlderWithSession) {
	a.Mux.HandleFunc(pattern, a.WithAuth(handler))
}

func (a App) HandleFuncWithNoAuth(pattern string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(pattern, a.WithNoAuth(handler))
}

func (a App) LoadRoutes() {
	// Auth
	a.HandleFuncWithNoAuth("POST /v1/auth/login", a.PasswordLoginHandler)
	a.HandleFuncWithNoAuth("POST /v1/auth/signup", a.PasswordSignupHandler)
	a.HandleFunc("GET /v1/auth/login/google", a.OAuthLoginHandler)
	a.HandleFunc("GET /v1/auth/callback/google", a.OAuthSignupHandler)
	a.HandleFuncWithAuth("GET /v1/auth/logout", a.LogoutSessionHandler)
	a.HandleFuncWithAuth("GET /v1/auth/logout_all", a.LogoutAllSessionsHandler)

	// Handle all other requests
	fs := http.FileServer(http.Dir("dist"))
	a.HandleFunc("GET /", a.StaticHandler("dist", "index.html", fs))
}
