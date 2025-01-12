package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
)

func (a *App) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Unable to upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.CloseNow()

	ctx := context.Background()

	for {
		_, msg, err := conn.Read(ctx)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			fmt.Printf("got err: %s", err)
			return
		}

		conn.Write(ctx, websocket.MessageText, msg)
	}
}
