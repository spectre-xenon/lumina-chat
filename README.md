# Lumina Chat

A group chat application built with Go and React, using websockets for real-time communication. This project is currently a work-in-progress but is nearing completion.

## Overview

Lumina Chat provides a real-time group chat experience. Key features include session-based authentication, chat creation and management, role-based actions, and a visually appealing interface.

**Key Technologies:** Go, React, PostgreSQL, Websockets

## Features

*   **Real-time Chat:** Websocket-based communication with a pub/sub pattern for instant messaging.
*   **Secure Authentication:** Session-based authentication and Argon2id hashing protect user data.
*   **Chat Management:** Create and manage chat groups with role-based permissions.
*   **Appealing Interface:** Modern design for a user-friendly experience.
*   **Pagination:** Infinite scroll for browsing chat messages.

## Technologies Used

*   **Backend:**
    *   Go
    *   SQLC (SQL Compiler): A code generation tool used to generate database boilerplate code from SQL queries.
    *   PostgreSQL
    *   Websockets (Implemented with a Pub/Sub pattern)
*   **Frontend:**
    *   React
    *   Tailwind CSS
*   **Infrastructure:**
    *   Docker

## Installation and Setup

**Prerequisites:** `go`, `pnpm`

1.  Clone the repository:
    ```bash
    git clone github.com/spectre-xenon/lumina-chat.git
    ```
2.  Navigate to the project directory:
    ```bash
    cd lumina-chat
    ```
3.  Install web dependencies and build static files:
    ```bash
    pnpm i && pnpm build
    ```
4.  Build the Go application:
    ```bash
    go build -v -o ./app ./cmd/server/main.go
    ```
5.  Create a `.env` file (or set environment variables) with the following:
    ```
    ORIGIN=http://localhost:8000
    DATABASE_URL=your_database_url
    GOOGLE_CLIENT_ID=your_google_client_id
    GOOGLE_CLIENT_SECRET=your_google_client_secret
    ```
6.  Run the application:
    ```bash
    ./app
    ```

## Development

**Prerequisites:** `go`, `pnpm`, `sqlc`, `air`, `docker`

1.  Run Postgres with Docker:
    ```bash
    docker compose up -d
    ```
2.  Run database migrations found in `Internal/db/migrations`. For example, using `golang-migrate`:
    ```bash
    migrate -database "$DATABASE_URL" -path Internal/db/migrations up
    ```
3.  Run the application with `air` for auto-rebuilding.
4.  Migrations can be added to the same `Internal/db/migrations` folder and rebuilt with SQLC.

---

Thank you for checking out Lumina Chat! Feel free to explore the code and contribute.
