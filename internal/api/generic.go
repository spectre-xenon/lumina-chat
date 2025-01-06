package api

import (
	"encoding/json"
	"net/http"
)

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
