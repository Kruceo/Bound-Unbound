package v1

import (
	"encoding/json"
	"net/http"
	"unbound-mngr-host/utils"
)

type Response[T any] struct {
	Message string
	Data    T
}

type ErrorResponse struct {
	Message   string
	ErrorCode string
	Error     bool
}

func FastErrorResponse(w http.ResponseWriter, r *http.Request, errorCode string, statusCode int) {
	b := ErrorResponse{Message: GetErrorMessage(errorCode), ErrorCode: errorCode, Error: true}
	encoded, err := json.Marshal(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(encoded)
}

func CorsHandler(w http.ResponseWriter, r *http.Request, methods string) bool {
	CORS_ORIGIN := utils.GetEnvOrDefault("CORS_ORIGIN", "*")
	w.Header().Set("Access-Control-Allow-Origin", CORS_ORIGIN)
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
		return true
	}
	return false
}
