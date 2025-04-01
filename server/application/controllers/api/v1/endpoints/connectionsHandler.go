package endpoints

import (
	"encoding/json"
	"net/http"

	v1 "server2/application/controllers/api/v1"
	usecases "server2/application/useCases"
)

func (bh V1Handlers) ConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}

	if _, err := v1.JWTMiddleware(w, r); err != nil {
		return
	}

	if r.Method != "GET" {
		return
	}

	getNode := usecases.GetNodeUseCase{Repo: &bh.NodeRepo}

	type ConnectionR struct {
		Name          string
		RemoteAddress string
	}
	var b v1.Response[[]ConnectionR] = v1.Response[[]ConnectionR]{Data: make([]ConnectionR, 0), Message: ""}
	for _, v := range bh.NodeRepo.IDs() {
		node := getNode.Execute(v)
		b.Data = append(b.Data, ConnectionR{Name: node.Name, RemoteAddress: v})
	}

	decoded, err := json.Marshal(b)
	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
