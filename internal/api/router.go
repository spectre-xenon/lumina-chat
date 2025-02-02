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
	apiRouter := http.NewServeMux()
	// Auth
	//  Password auth
	a.handleFuncWithNoAuth(apiRouter, "POST /auth/login", a.passwordLoginHandler)
	a.handleFuncWithNoAuth(apiRouter, "POST /auth/signup", a.passwordSignupHandler)
	//  OAuth auth
	a.handleFunc(apiRouter, "GET /auth/login/google", a.oauthLoginHandler)
	a.handleFuncWithNoAuth(apiRouter, "GET /auth/callback/google", a.oauthSignupHandler)
	//  Logging out
	a.handleFuncWithAuth(apiRouter, "GET /auth/logout", a.logoutSessionHandler)
	a.handleFuncWithAuth(apiRouter, "GET /auth/logout_all", a.logoutAllSessionsHandler)
	//  Check auth status
	a.handleFuncWithAuth(apiRouter, "GET /auth", func(w http.ResponseWriter, r *http.Request, session db.Session) {
		w.WriteHeader(http.StatusOK)
	})

	// User
	a.handleFuncWithAuth(apiRouter, "PATCH /user", a.patchUserHandler)

	// Chat
	a.handleFuncWithAuth(apiRouter, "GET /chats", a.getUserChatsHandler)

	// Enable logging on dev enviroments
	env := os.Getenv("LUMINA_ENV")
	var handlerStack middleware.Middleware
	if env != "prod" {
		handlerStack = middleware.CreateStack(
			middleware.Logging,
			middleware.LimitBodySize,
		)
	} else {
		handlerStack = middleware.CreateStack(
			middleware.LimitBodySize,
		)
	}

	// Handle api routes
	a.mux.Handle("/v1/", http.StripPrefix("/v1", handlerStack(apiRouter)))

	// Handle all other static requests
	fs := http.FileServer(http.Dir("dist"))
	a.mux.HandleFunc("/", a.staticHandler("dist", "index.html", fs))
}
