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

type RedirectBody struct {
	From       string
	RecordType string
	To         string
	LocalZone  bool
}

func RedirectAddressHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, POST, DELETE, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method == "GET" {

		connectionName := r.PathValue("connection")
		client, exists := memory.Connections[connectionName]
		if !exists {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}
		id := fmt.Sprintf("%x", rand.Int())
		client.Send(id+" list redirects", true)

		memory.WaitForResponse(id)

		var b v1.Response[[]RedirectBody] = v1.Response[[]RedirectBody]{Data: make([]RedirectBody, 0)}
		for _, v := range strings.Split(memory.Responses[id], ",") {
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
		client.Send(fmt.Sprintf("%s redirect %s %s %s %s", id, b.From, b.RecordType, b.To, localzoneStr), true)

		memory.WaitForResponse(id)
		w.Write(nil)
		return
	} else if r.Method == "DELETE" {
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

		var b struct{ Domain string }
		err = json.Unmarshal(body, &b)

		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
			return
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		client.Send(id+" unredirect "+b.Domain, true)

		memory.WaitForResponse(id)

		w.Write(nil)
		return
	}
}
