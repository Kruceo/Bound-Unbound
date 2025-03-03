package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/commands"

	"github.com/gorilla/websocket"
)

func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}

	if r.Method != "POST" {
		return
	}

	connectionName := r.PathValue("connection")
	conn := commands.Connections[connectionName]
	if conn == nil {
		fmt.Println("Not found:", connectionName)
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
		return
	}
	id := fmt.Sprintf("%x", rand.Int())
	conn.WriteMessage(websocket.TextMessage, []byte(id+" reload"))

	commands.WaitForResponse(id)

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
