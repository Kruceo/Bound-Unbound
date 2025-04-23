package v1

import (
	"encoding/json"
	"net/http"

	"server2/application/presentation"
)

type ConnectionW struct {
	Name          string `json:"name"`
	RemoteAddress string `json:"remoteAddress"`
}

func (bh V1APIHandlers) ConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	var b presentation.Response[[]ConnectionW] = presentation.Response[[]ConnectionW]{Data: make([]ConnectionW, 0), Message: ""}
	for _, v := range bh.nodePersistenceUseCase.IDs() {
		node := bh.nodePersistenceUseCase.Get(v)
		b.Data = append(b.Data, ConnectionW{Name: node.Name, RemoteAddress: v})
	}

	decoded, err := json.Marshal(b)
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
