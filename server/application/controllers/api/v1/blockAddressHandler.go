package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"server2/application/presentation"
	usecases "server2/application/useCases"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var cipherMessage = usecases.CipherCommandMessageUseCase{}

type BlockedNamesW struct {
	Names []string `json:"names"`
}
type BlockedNamesR struct {
	Names []string `json:"names"`
}

func (bh *V1APIHandlers) BlockAddressHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	connectionName := vars["connection"]

	if r.Method == "GET" {

		client, err := bh.nodePersistenceUseCase.Get(connectionName)
		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		id := fmt.Sprintf("%x", rand.Int())

		encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s list blocked", id), client.Cipher)

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
		var b presentation.Response[BlockedNamesW]

		rawNames, err := bh.responseRepo.ReadResponse(id)
		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}

		b.Data.Names = strings.Split(rawNames, ",")

		if len(b.Data.Names) == 1 && b.Data.Names[0] == "" {
			b.Data.Names = []string{}
		}
		encoded, err := json.Marshal(b)
		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(encoded)
		// }()
		return
	} else if r.Method == "POST" {

		client, err := bh.nodePersistenceUseCase.Get(connectionName)
		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		var body []byte
		body, err = io.ReadAll(r.Body)
		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
			return
		}

		var b BlockedNamesR
		err = json.Unmarshal(body, &b)

		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusBadRequest)
			return
		}

		if len(b.Names) < 1 {
			bh.fastErrorResponses.Execute(w, r, "BODY_FORMAT", http.StatusBadRequest)
			return
		}

		id := fmt.Sprintf("%X", rand.Int())
		// client.Send(, true)

		encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s block %s", id, strings.Join(b.Names, ",")), client.Cipher)

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
		bh.responseRepo.ReadResponse(id)
		w.Write(nil)
		return
	} else if r.Method == "DELETE" {
		client, err := bh.nodePersistenceUseCase.Get(connectionName)
		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		var body []byte
		body, err = io.ReadAll(r.Body)
		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
			return
		}

		var b BlockedNamesR
		err = json.Unmarshal(body, &b)

		if err != nil {
			bh.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
			return
		}

		id := fmt.Sprintf("%X", rand.Int())
		encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s unblock %s", id, strings.Join(b.Names, ",")), client.Cipher)

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
		w.Write(nil)
		return
	}
}
