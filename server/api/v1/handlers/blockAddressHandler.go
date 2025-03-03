package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/commands"

	"github.com/gorilla/websocket"
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
		conn := commands.Connections[connectionName]
		if conn == nil {
			fmt.Println("Not found:", connectionName)
			w.WriteHeader(http.StatusNotFound)
			w.Write(nil)
			return
		}
		id := fmt.Sprintf("%x", rand.Int())
		conn.WriteMessage(websocket.TextMessage, []byte(id+" list blocked"))

		commands.WaitForResponse(id)

		var b v1.Response[BlockedNames]
		b.Data.Names = strings.Split(commands.Responses[id], ",")
		if len(b.Data.Names) == 1 && b.Data.Names[0] == "" {
			b.Data.Names = []string{}
		}
		decoded, err := json.Marshal(b)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte{})
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(decoded)
		return
	} else if r.Method == "POST" {

		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")
		conn := commands.Connections[connectionName]
		if conn == nil {
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
		}

		var b BlockedNames
		err = json.Unmarshal(body, &b)

		if err != nil {
			fmt.Println("json decode error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		conn.WriteMessage(websocket.TextMessage, []byte(id+" block "+strings.Join(b.Names, ",")))

		commands.WaitForResponse(id)

		w.Write(nil)
		return
	} else if r.Method == "DELETE" {

		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")
		conn := commands.Connections[connectionName]
		if conn == nil {
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
		}

		var b BlockedNames
		err = json.Unmarshal(body, &b)

		if err != nil {
			fmt.Println("json decode error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		conn.WriteMessage(websocket.TextMessage, []byte(id+" unblock "+strings.Join(b.Names, ",")))

		commands.WaitForResponse(id)

		w.Write(nil)
		return
	}
}
