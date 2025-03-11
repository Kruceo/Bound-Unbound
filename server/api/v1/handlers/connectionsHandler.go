package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/memory"
)

func ConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}

	if r.Method != "GET" {
		return
	}
	type ConnectionR struct {
		Name          string
		RemoteAddress string
	}
	var b v1.Response[[]ConnectionR] = v1.Response[[]ConnectionR]{Data: make([]ConnectionR, 0), Message: ""}
	for k, v := range memory.Connections {
		b.Data = append(b.Data, ConnectionR{Name: v.Name, RemoteAddress: k})
	}

	decoded, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
