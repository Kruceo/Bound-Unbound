package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	v1 "server2/api/v1"
	usecases "server2/application/useCases"
	"strings"

	"github.com/gorilla/websocket"
)

type RedirectBody struct {
	From       string
	RecordType string
	To         string
	LocalZone  bool
}

func (bh *V1Handlers) RedirectAddressHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, POST, DELETE, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	getNode := usecases.GetNodeUseCase{Repo: bh.NodeRepo}

	if r.Method == "GET" {

		connectionName := r.PathValue("connection")
		client := getNode.Execute(connectionName)
		if client == nil {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		id := fmt.Sprintf("%x", rand.Int())
		encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s list redirects", id), &client.Cipher)

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

		var b v1.Response[[]RedirectBody] = v1.Response[[]RedirectBody]{Data: make([]RedirectBody, 0)}

		rawRedirects, err := bh.ResponseRepo.ReadResponse(id)
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}

		for _, v := range strings.Split(rawRedirects, ",") {
			if v == "" {
				continue
			}
			vsplt := strings.Split(v, " ")
			from := vsplt[0]
			rtype := vsplt[1]
			to := vsplt[2]
			localZone := vsplt[3] == "true"
			b.Data = append(b.Data, RedirectBody{From: from, RecordType: rtype, To: to, LocalZone: localZone})
		}

		decoded, err := json.Marshal(b)
		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_ENCODING", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(decoded)
		return
	} else if r.Method == "POST" {
		connectionName := r.PathValue("connection")
		client := getNode.Execute(connectionName)
		if client == nil {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			v1.FastErrorResponse(w, r, "READ_BODY", http.StatusInternalServerError)
			return
		}

		var b RedirectBody
		err = json.Unmarshal(body, &b)

		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
			return
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		localzoneStr := ""
		if b.LocalZone {
			localzoneStr = "local-zone"
		}
		// client.Send(/, true)

		encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s redirect %s %s %s %s", id, b.From, b.RecordType, b.To, localzoneStr), &client.Cipher)

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
		w.Write(nil)
		return
	} else if r.Method == "DELETE" {
		connectionName := r.PathValue("connection")
		client := getNode.Execute(connectionName)
		if client == nil {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			v1.FastErrorResponse(w, r, "READ_BODY", http.StatusInternalServerError)
			return
		}

		var b struct{ Domain string }
		err = json.Unmarshal(body, &b)

		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
			return
		}

		id := fmt.Sprintf("%X", rand.Int())
		encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("%s unredirect %s", id, b.Domain), &client.Cipher)

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

		w.Write(nil)
		return
	}
}
