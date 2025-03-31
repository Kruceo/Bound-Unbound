package endpoints

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	v1 "server2/application/controllers/api/v1"
	usecases "server2/application/useCases"

	"github.com/gorilla/websocket"
)

func (bh *V1Handlers) ConfigHashHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method != "GET" {
		return
	}

	getNode := usecases.GetNodeUseCase{Repo: bh.NodeRepo}

	type HashR struct {
		Hash string
	}

	connectionName := r.PathValue("connection")
	client := getNode.Execute(connectionName)
	if client == nil {
		v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
		return
	}

	id := fmt.Sprintf("%x", rand.Int())
	encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s list confighash", id), &client.Cipher)

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

	res, err := bh.ResponseRepo.ReadResponse(id)
	if err != nil {
		v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	var b v1.Response[HashR] = v1.Response[HashR]{Data: HashR{Hash: res}, Message: ""}

	decoded, err := json.Marshal(b)
	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
