package api

import (
	"net/http"

	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/middleware"
)

type App struct {
	db  *db.Queries
	Mux *http.ServeMux
}

func New(db *db.Queries, mux *http.ServeMux) App {
	return App{db, mux}
}

func (a *App) LoadRoutes() {
	// Auth
	a.Mux.HandleFunc("POST /v1/auth/login", a.PasswordLoginHandler)
	a.Mux.HandleFunc("POST /v1/auth/signup", a.PasswordSignupHandler)
	a.Mux.HandleFunc("GET /v1/auth/login/google", a.OAuthLoginHandler)
	a.Mux.HandleFunc("GET /v1/auth/callback/google", a.OAuthSignupHandler)
	a.Mux.HandleFunc("GET /v1/auth/logout", a.LogoutSessionHandler)
	a.Mux.HandleFunc("GET /v1/auth/logout_all", a.LogoutAllSessionsHandler)

	// Handle all other requests
	fs := http.FileServer(http.Dir("dist"))
	a.Mux.HandleFunc("GET /", middleware.StaticHandler("dist", "index.html", fs))
}
