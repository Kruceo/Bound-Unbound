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

	if r.Method == "GET" {

		connectionName := r.PathValue("connection")
		conn := commands.Connections[connectionName]
		if conn == nil {
			fmt.Println("Not found:", connectionName)
			w.WriteHeader(http.StatusNotFound)
			w.Write(nil)
			return
		}
		id := fmt.Sprintf("%x", rand.Int())
		conn.WriteMessage(websocket.TextMessage, []byte(id+" list redirects"))

		commands.WaitForResponse(id)

		var b v1.Response[[]RedirectBody] = v1.Response[[]RedirectBody]{Data: make([]RedirectBody, 0)}
		for _, v := range strings.Split(commands.Responses[id], ",") {
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
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte{})
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(decoded)
		return
	} else if r.Method == "POST" {
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

		var b RedirectBody
		err = json.Unmarshal(body, &b)

		if err != nil {
			fmt.Println("json decode error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s redirect %s %s %s %v", id, b.From, b.RecordType, b.To, b.LocalZone)))

		commands.WaitForResponse(id)
		w.Write(nil)
		return
	} else if r.Method == "DELETE" {
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

		var b struct{ Domain string }
		err = json.Unmarshal(body, &b)

		if err != nil {
			fmt.Println("json decode error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}

		id := fmt.Sprintf("%X", rand.Int()*1000)
		conn.WriteMessage(websocket.TextMessage, []byte(id+" unredirect "+b.Domain))

		commands.WaitForResponse(id)

		w.Write(nil)
		return
	}
}
