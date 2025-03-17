package v1

import (
	"net/http"
	"unbound-mngr-host/utils"
)

type Response[T any] struct {
	Message string
	Data    T
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
