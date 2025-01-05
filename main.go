package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/spectre-xenon/lumina-chat/internal/middleware"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()

	// Handle all other requests
	fs := http.FileServer(http.Dir("dist"))
	mux.HandleFunc("GET /", middleware.StaticHandler("dist", "index.html", fs))

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
