package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/spectre-xenon/lumina-chat/internal/middleware"
)

const STATIC_PATH = "dist"
const INDEX_PATH = "index.html"

func main() {
	godotenv.Load()

	mux := http.NewServeMux()
	m := middleware.Middleware{
		StaticPath: STATIC_PATH,
		IndexPath:  INDEX_PATH,
	}

	static := http.FileServer(http.Dir(STATIC_PATH))
	mux.HandleFunc("/", m.Root(static))

	// host:port
	addr := "127.0.0.1:8000"
	server := &http.Server{
		Handler: mux,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server is listening on: %s \n", addr)
	server.ListenAndServe()
}
