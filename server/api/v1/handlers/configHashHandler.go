package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/commands"

	"github.com/gorilla/websocket"
)

func ConfigHashHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}

	if r.Method != "GET" {
		return
	}
	type HashR struct {
		Hash string
	}

	connectionName := r.PathValue("connection")
	conn := commands.Connections[connectionName]
	id := fmt.Sprintf("%x", rand.Int())
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s list confighash", id)))

	commands.WaitForResponse(id)
	res := commands.Responses[id]

	var b v1.Response[HashR] = v1.Response[HashR]{Data: HashR{Hash: res}, Message: ""}

	decoded, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
