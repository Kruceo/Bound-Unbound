package v1

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"server2/application/presentation"
	usecases "server2/application/useCases"

	"github.com/gorilla/websocket"
)

func (bh *V1APIHandlers) ConfigHashHandler(w http.ResponseWriter, r *http.Request) {
	getNode := usecases.GetNodeUseCase{Repo: &bh.nodeRepo}

	type HashR struct {
		Hash string
	}

	connectionName := r.PathValue("connection")
	client := getNode.Execute(connectionName)
	if client == nil {
		bh.fastErrorResponses.Execute(w, r, "UNKNOWN_NODE", http.StatusNotFound)
		return
	}

	id := fmt.Sprintf("%x", rand.Int())
	encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s list confighash", id), &client.Cipher)

	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "CONNECTION_SECURITY", http.StatusInternalServerError)
		return
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, encryptedMessage)
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	err = bh.responseRepo.WaitForResponse(id)

	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	res, err := (bh.responseRepo).ReadResponse(id)
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	var b presentation.Response[HashR] = presentation.Response[HashR]{Data: HashR{Hash: res}, Message: ""}

	decoded, err := json.Marshal(b)
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
