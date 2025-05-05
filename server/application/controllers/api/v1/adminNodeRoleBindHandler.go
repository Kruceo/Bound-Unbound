package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server2/application/presentation"
)

type getBindsNodeW struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getBindsRoleW struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type GetBindsW struct {
	ID   string        `json:"id"`
	Node getBindsNodeW `json:"node"`
	Role getBindsRoleW `json:"role"`
}

func (ad *V1AdminHandlers) AdminGetBinds(w http.ResponseWriter, r *http.Request) {
	binds, err := ad.bindsPersistence.GetAllWithIncluded(2048)
	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "ROLE_REPO", http.StatusInternalServerError)
		return
	}

	var response presentation.Response[[]GetBindsW]

	for _, v := range binds {
		fmt.Println("id", v.ID)
		response.Data = append(response.Data, GetBindsW{
			ID: v.ID,
			Node: getBindsNodeW{
				ID:   v.Node.ID,
				Name: v.Node.Name,
			},
			Role: getBindsRoleW{
				ID:          v.Role.ID,
				Name:        v.Role.Name,
				Permissions: v.Role.Permissions,
			},
		})
	}

	encoded, err := json.Marshal(response)

	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(encoded)
}

type BindPostR struct {
	NodeID string `json:"nodeId"`
	RoleID string `json:"roleId"`
}

type BindPostW struct {
	ID string `json:"id"`
}

func (ad *V1AdminHandlers) AdminPostBinds(w http.ResponseWriter, r *http.Request) {
	var response presentation.Response[[]PostRoleW]

	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b []BindPostR
	err = json.Unmarshal(body, &b)

	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	for _, v := range b {
		generatedID, err := ad.bindsPersistence.Save(v.NodeID, v.RoleID)
		if err != nil {
			ad.fastErrorResponses.Execute(w, r, "BIND_SAVE", http.StatusInternalServerError)
			return
		}
		response.Data = append(response.Data, PostRoleW{ID: generatedID})
	}

	encoded, err := json.Marshal(response)

	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(encoded)
}

type BindDeleteR struct {
	ID string `json:"id"`
}

func (ad *V1AdminHandlers) AdminDeleteBinds(w http.ResponseWriter, r *http.Request) {
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b []BindDeleteR
	err = json.Unmarshal(body, &b)

	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	for _, v := range b {
		err := ad.bindsPersistence.Delete(v.ID)
		if err != nil {
			ad.fastErrorResponses.Execute(w, r, "BIND_DELETION", http.StatusInternalServerError)
			return
		}
	}
	var response presentation.Response[bool]
	response.Data = true
	encoded, err := json.Marshal(response)

	if err != nil {
		ad.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(encoded)
}
