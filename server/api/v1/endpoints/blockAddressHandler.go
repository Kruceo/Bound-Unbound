package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	// "server/bh.ResponseRepo"
	v1 "server2/api/v1"
	"server2/application/entities"
	usecases "server2/application/useCases"
	"strings"
)

type V1Handlers struct {
	NodeRepo     entities.NodeRepository
	ResponseRepo entities.ResponsesReporisory
}

func (bh *V1Handlers) BlockAddressHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, POST, DELETE, OPTIONS") {
		return
	}

	getNode := usecases.GetNodeUseCase{Repo: bh.NodeRepo}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method == "GET" {
		type BlockedNames struct {
			Names []string
		}
		connectionName := r.PathValue("connection")

		client := getNode.Execute(connectionName)
		if client == nil {
			v1.FastErrorResponse(w, r, "UNKNOWN_NODE", http.StatusNotFound)
			return
		}

		id := fmt.Sprintf("%x", rand.Int())
		client.Send(id+" list blocked", true)

		fmt.Println("Waiting for bh")
		err := bh.ResponseRepo.WaitForResponse(id)
		fmt.Println("pass")
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}
		var b v1.Response[BlockedNames]

		rawNames, err := bh.ResponseRepo.ReadResponse(id)
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}

		b.Data.Names = strings.Split(rawNames, ",")

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

		err = bh.ResponseRepo.WaitForResponse(id)
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}
		bh.ResponseRepo.ReadResponse(id)
		w.Write(nil)
		return
	} else if r.Method == "DELETE" {

		type BlockedNames struct {
			Names []string
		}
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

		var b BlockedNames
		err = json.Unmarshal(body, &b)

		if err != nil {
			v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
			return
		}

		id := fmt.Sprintf("%X", rand.Int())
		client.Send(id+" unblock "+strings.Join(b.Names, ","), true)

		err = bh.ResponseRepo.WaitForResponse(id)
		if err != nil {
			v1.FastErrorResponse(w, r, "NODE_RESPONSE", http.StatusInternalServerError)
			return
		}
		w.Write(nil)
		return
	}
}
