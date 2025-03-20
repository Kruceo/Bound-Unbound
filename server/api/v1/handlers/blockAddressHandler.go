package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/memory"
)

func BlockAddressHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, POST, DELETE, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method == "GET" {
		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")
		client, exists := memory.Connections[connectionName]
		if !exists {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		id := fmt.Sprintf("%x", rand.Int())
		client.Send(id+" list blocked", true)

		err := memory.WaitForResponse(id)
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}
		var b v1.Response[BlockedNames]

		b.Data.Names = strings.Split(memory.ReadResponse(id), ",")

		if len(b.Data.Names) == 1 && b.Data.Names[0] == "" {
			b.Data.Names = []string{}
		}
		encoded, err := json.Marshal(b)
		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_ENCODING", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(encoded)
		// }()
		return
	} else if r.Method == "POST" {

		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")
		client, exists := memory.Connections[connectionName]
		if !exists {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			v1.FastErrorResponse(w, r, "READ_BODY", http.StatusInternalServerError)
			return
		}

		var b BlockedNames
		err = json.Unmarshal(body, &b)

		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusBadRequest)
			return
		}

		if len(b.Names) < 1 {
			v1.FastErrorResponse(w, r, "BODY_FORMAT", http.StatusBadRequest)
			return
		}

		id := fmt.Sprintf("%X", rand.Int())
		client.Send(id+" block "+strings.Join(b.Names, ","), true)

		err = memory.WaitForResponse(id)
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}
		memory.ReadResponse(id)
		w.Write(nil)
		return
	} else if r.Method == "DELETE" {

		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")
		client, exists := memory.Connections[connectionName]
		if !exists {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			v1.FastErrorResponse(w, r, "READ_BODY", http.StatusInternalServerError)
			return
		}

		var b BlockedNames
		err = json.Unmarshal(body, &b)

		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
			return
		}

		id := fmt.Sprintf("%X", rand.Int())
		client.Send(id+" unblock "+strings.Join(b.Names, ","), true)

		err = memory.WaitForResponse(id)
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}
		memory.ReadResponse(id)
		w.Write(nil)
		return
	}
}
