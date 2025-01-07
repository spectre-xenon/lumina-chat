package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spectre-xenon/lumina-chat/internal/api"
	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/middleware"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	// New connection pool
	conn, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panicf("Error connecting to database: %s", err)
	}
	defer conn.Close()

	db := db.New(conn)

	// Create new api instance and define all routes
	app := api.New(db, http.NewServeMux())
	app.LoadRoutes()

	// host:port
	addr := "127.0.0.1:8000"
	server := &http.Server{
		Handler: middleware.Logging(app.Mux),
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server is listening on: %s \n", os.Getenv("ORIGIN"))
	server.ListenAndServe()
}
