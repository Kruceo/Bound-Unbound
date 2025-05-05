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

	// gets the reader defined at specific role middleware
	nodes, err := bh.nodeRoleBindUseCase.GetNodesForRole(r.Header.Get("X-Role-ID"))
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "RECOVER_NODES", http.StatusInternalServerError)
		return
	}
	var b presentation.Response[[]ConnectionW] = presentation.Response[[]ConnectionW]{Data: make([]ConnectionW, 0), Message: ""}
	for _, v := range nodes {
		b.Data = append(b.Data, ConnectionW{Name: v.Name, RemoteAddress: v.ID})
	}
	decoded, err := json.Marshal(b)
	if err != nil {
		bh.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(decoded)

}
