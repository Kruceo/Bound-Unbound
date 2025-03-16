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

	if r.Method == "GET" {
		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")
		client, exists := memory.Connections[connectionName]
		if !exists {
			fmt.Println("Not found:", connectionName)
			w.WriteHeader(http.StatusNotFound)
			w.Write(nil)
			return
		}
		id := fmt.Sprintf("%x", rand.Int())
		client.Send(id+" list blocked", true)

		memory.WaitForResponse(id)

		var b v1.Response[BlockedNames]
		b.Data.Names = strings.Split(memory.ReadResponse(id), ",")
		if len(b.Data.Names) == 1 && b.Data.Names[0] == "" {
			b.Data.Names = []string{}
		}
		decoded, err := json.Marshal(b)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte{})
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(decoded)
		return
	} else if r.Method == "POST" {

		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")
		client, exists := memory.Connections[connectionName]
		if !exists {
			fmt.Println("Not found:", connectionName)
			w.WriteHeader(http.StatusNotFound)
			w.Write(nil)
			return
		}

		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("body read error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}

		var b BlockedNames
		err = json.Unmarshal(body, &b)

		if err != nil {
			fmt.Println("json decode error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		client.Send(id+" block "+strings.Join(b.Names, ","), true)

		memory.WaitForResponse(id)
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
			fmt.Println("Not found:", connectionName)
			w.WriteHeader(http.StatusNotFound)
			w.Write(nil)
			return
		}

		var body []byte
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("body read error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}

		var b BlockedNames
		err = json.Unmarshal(body, &b)

		if err != nil {
			fmt.Println("json decode error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		client.Send(id+" unblock "+strings.Join(b.Names, ","), true)

		memory.WaitForResponse(id)
		memory.ReadResponse(id)
		w.Write(nil)
		return
	}
}
