package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/util"
)

var mimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func (a *App) getUserChatsHandler(w http.ResponseWriter, r *http.Request, session db.Session) {
	chats, err := a.db.GetUserChats(r.Context(), session.UserID)
	if err != nil && errors.Is(pgx.ErrNoRows, err) == false {
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	response := ApiResponse[db.GetUserChatsRow]{Data: chats}
	JSON(w, response)
}

func (a *App) createChatHandler(w http.ResponseWriter, r *http.Request, session db.Session) {
	// Parse multipart data
	if err := r.ParseMultipartForm(0); err != nil {
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	name, exists := r.Form["name"]
	if !exists || len(name[0]) < 3 {
		http.Error(w, "Name too short", http.StatusBadRequest)
		return
	}

	// Handle picture
	file, header, err := r.FormFile("picture")
	var pictureUrl *string
	if err == nil {
		if header.Size > 1<<20 || mimeTypes[header.Header.Get("Content-Type")] == false {
			http.Error(w, "File bigger than 1mb or unsupportted file type", http.StatusBadRequest)
			return
		}

		saveUrl, err := util.SavePicture(file, *header)
		if err != nil {
			log.Printf("Failed to save file: %s\n", err)
			internalServerError(w)
			return
		}

		pictureUrl = &saveUrl
		file.Close()
	}

	ctx := context.Background()
	// Save chat to db
	inviteLink := "/invite/" + uuid.New().String()
	chat, err := a.db.CreateChat(ctx, db.CreateChatParams{
		Name:       name[0],
		Picture:    pictureUrl,
		InviteLink: &inviteLink,
	})
	if err != nil {
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	// Add user as member
	err = a.db.AddChatMember(ctx, db.AddChatMemberParams{ChatID: chat.ID, UserID: session.UserID})
	if err != nil {
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	msg, err := a.db.CreateMessage(ctx, db.CreateMessageParams{
		ChatID:   chat.ID,
		SenderID: session.UserID,
		Content:  "Chat init message!",
	})
	if err != nil {
		log.Printf("Database error: %s\n", err)
		internalServerError(w)
		return
	}

	chatWithMsg := db.GetUserChatsRow{
		ID:         chat.ID,
		Name:       chat.Name,
		InviteLink: chat.InviteLink,
		Picture:    chat.Picture,
		Message:    msg,
	}

	response := ApiResponse[db.GetUserChatsRow]{Data: []db.GetUserChatsRow{chatWithMsg}}
	JSON(w, response)
}
