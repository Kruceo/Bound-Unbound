package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/memory"
)

func ConfigHashHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method != "GET" {
		return
	}
	type HashR struct {
		Hash string
	}

	connectionName := r.PathValue("connection")
	client, exists := memory.Connections[connectionName]
	if !exists {
		v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
		return
	}
	id := fmt.Sprintf("%x", rand.Int())
	client.Send(fmt.Sprintf("%s list confighash", id), true)

	memory.WaitForResponse(id)
	res := memory.ReadResponse(id)

	var b v1.Response[HashR] = v1.Response[HashR]{Data: HashR{Hash: res}, Message: ""}

	decoded, err := json.Marshal(b)
	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
