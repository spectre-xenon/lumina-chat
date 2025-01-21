package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spectre-xenon/lumina-chat/internal/api"
	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/workerpool"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	// New connection pool
	conn, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	_, err2 := conn.Query(ctx, "select username from users limit 1;")
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	if err2 != nil {
		log.Fatalf("Error connection test: %s\n", err2)
	}
	defer conn.Close()

	db := db.New(conn)

	workerPool := workerpool.NewWorkerPool(runtime.NumCPU())
	defer workerPool.Close()

	// Create new api instance and define all routes
	app := api.New(db, http.NewServeMux(), workerPool)
	app.LoadRoutes()

	// host:port
	addr := "127.0.0.1:8000"
	server := &http.Server{
		Handler: app.GetHandler(),
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server is listening on: %s \n", os.Getenv("ORIGIN"))
	server.ListenAndServe()
}
