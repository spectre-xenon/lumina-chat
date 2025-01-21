package api

import (
	"net/http"
	"os"

	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/middleware"
	"github.com/spectre-xenon/lumina-chat/internal/workerpool"
)

type App struct {
	db         *db.Queries
	mux        *http.ServeMux
	workerPool *workerpool.WorkerPool
}

func New(db *db.Queries, mux *http.ServeMux, workerPool *workerpool.WorkerPool) App {
	return App{db, mux, workerPool}
}

func (a *App) handleFunc(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.HandleFunc(pattern, handler)
}

func (a *App) handleFuncWithAuth(mux *http.ServeMux, pattern string, handler HanlderWithSession) {
	mux.HandleFunc(pattern, a.WithAuth(handler))
}

// Requires there to be no auth
func (a *App) handleFuncWithNoAuth(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.HandleFunc(pattern, a.WithNoAuth(handler))
}

func (a *App) GetHandler() *http.ServeMux {
	return a.mux
}

func (a *App) LoadRoutes() {
	loggedRouter := http.NewServeMux()

	// Auth
	//  Password auth
	a.handleFuncWithNoAuth(loggedRouter, "POST /v1/auth/login", a.PasswordLoginHandler)
	a.handleFuncWithNoAuth(loggedRouter, "POST /v1/auth/signup", a.PasswordSignupHandler)
	//  OAuth auth
	a.handleFunc(loggedRouter, "GET /v1/auth/login/google", a.OAuthLoginHandler)
	a.handleFuncWithNoAuth(loggedRouter, "GET /v1/auth/callback/google", a.OAuthSignupHandler)
	//  Logging out
	a.handleFuncWithAuth(loggedRouter, "GET /v1/auth/logout", a.LogoutSessionHandler)
	a.handleFuncWithAuth(loggedRouter, "GET /v1/auth/logout_all", a.LogoutAllSessionsHandler)
	//  Check auth status
	a.handleFuncWithAuth(loggedRouter, "GET /v1/auth", func(w http.ResponseWriter, r *http.Request, session db.Session) {
		w.WriteHeader(http.StatusOK)
	})

	// Handle all other static requests
	fs := http.FileServer(http.Dir("dist"))
	loggedRouter.HandleFunc("GET /", a.StaticHandler("dist", "index.html", fs))

	// Enable logging on dev enviroments
	env := os.Getenv("LUMINA_ENV")
	if env != "prod" {
		a.mux.Handle("/", middleware.Logging(loggedRouter))
	} else {

		a.mux.Handle("/", loggedRouter)
	}
}
