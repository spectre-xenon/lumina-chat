package api

import (
	"encoding/json"
	"net/http"

	"github.com/spectre-xenon/lumina-chat/internal/db"
	"github.com/spectre-xenon/lumina-chat/internal/util"
)

type HanlderWithSession func(w http.ResponseWriter, r *http.Request, session db.Session)

type ApiResponse[T any] struct {
	Data    []T  `json:"data"`
	ErrCode *int `json:"err_code"`
}

func JSONError(w http.ResponseWriter, err any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func JSON(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func internalServerError(w http.ResponseWriter) {
	response := ApiResponse[db.User]{ErrCode: util.Of(InternalServerError)}
	JSONError(w, response, http.StatusInternalServerError)
}
