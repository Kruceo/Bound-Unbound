package endpoints

import (
	"fmt"
	"math/rand"
	"net/http"

	v1 "server2/application/controllers/api/v1"
	usecases "server2/application/useCases"

	"github.com/gorilla/websocket"
)

func (bh *V1Handlers) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method != "POST" {
		return
	}

	getNode := usecases.GetNodeUseCase{Repo: bh.NodeRepo}

	connectionName := r.PathValue("connection")
	client := getNode.Execute(connectionName)
	if client == nil {
		v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
		return
	}
	id := fmt.Sprintf("%x", rand.Int())
	encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s reload", id), &client.Cipher)

	if err != nil {
		v1.FastErrorResponse(w, r, "CONNECTION_SECURITY", http.StatusInternalServerError)
		return
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, encryptedMessage)
	if err != nil {
		v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	err = bh.ResponseRepo.WaitForResponse(id)
	if err != nil {
		v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
