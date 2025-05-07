package v1

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func (bh *V1APIHandlers) ReloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		return
	}

	vars := mux.Vars(r)
	connectionName := vars["connection"]

	client, err := bh.nodePersistenceUseCase.Get(connectionName)
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "UNKNOWN_NODE", http.StatusNotFound)
		return
	}
	id := fmt.Sprintf("%x", rand.Int())

	encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s reload", id), client.Cipher)

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
	fmt.Println(err)
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
