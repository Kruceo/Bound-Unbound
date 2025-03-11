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

	if r.Method != "POST" {
		return
	}

	connectionName := r.PathValue("connection")
	client, exists := memory.Connections[connectionName]
	if !exists {
		fmt.Println("Not found:", connectionName)
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
		return
	}
	id := fmt.Sprintf("%x", rand.Int())
	client.Send(id+" reload", true)

	memory.WaitForResponse(id)

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
