package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/memory"
)

func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method != "POST" {
		return
	}

	connectionName := r.PathValue("connection")
	client, exists := memory.Connections[connectionName]
	if !exists {
		v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
		return
	}
	id := fmt.Sprintf("%x", rand.Int())
	client.Send(id+" reload", true)

	err := memory.WaitForResponse(id)
	if err != nil {
		v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
