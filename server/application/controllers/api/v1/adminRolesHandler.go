package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server2/application/presentation"
)

type GetRolesW struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (a *V1AdminHandlers) AdminGetRoles(w http.ResponseWriter, r *http.Request) {

	roles, err := a.roleUseCase.GetAll(258)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "ROLE_REPO", http.StatusInternalServerError)
	}

	var response presentation.Response[[]GetRolesW]

	for _, v := range roles {
		response.Data = append(response.Data, GetRolesW{Name: v.Name, ID: v.ID})
	}

	encoded, err := json.Marshal(response)

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(encoded)

}

type PostRoleR struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type PostRoleW struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (a *V1AdminHandlers) AdminPostRole(w http.ResponseWriter, r *http.Request) {

	var payload []PostRoleR
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var response presentation.Response[[]PostRoleW]

	for _, v := range payload {
		createdID, err := a.roleUseCase.Save(v.Name, v.Permissions)
		if err != nil {
			a.fastErrorResponses.Execute(w, r, "ROLE_CREATION", http.StatusBadRequest)
			return
		}

		response.Data = append(response.Data, PostRoleW{ID: createdID, Name: v.Name, Permissions: v.Permissions})
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(encoded)
}

type DeleteRoleW struct {
	ID string `json:"id"`
}

func (a *V1AdminHandlers) AdminDeleteRole(w http.ResponseWriter, r *http.Request) {

	var payload []DeleteRoleW
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("missing 'validation over existent users using role'")

	for _, v := range payload {
		if v.ID == "" {
			a.fastErrorResponses.Execute(w, r, "MISSING_ID", http.StatusBadRequest)
			return
		}

		usersWithThisRole, err := a.userUseCase.SearchByRoleID(v.ID)
		if err != nil {
			a.fastErrorResponses.Execute(w, r, "GET_ROLE_USERS", http.StatusInternalServerError)
			return
		}

		if len(usersWithThisRole) > 0 {
			a.fastErrorResponses.Execute(w, r, "USER_USING_ROLE", http.StatusInternalServerError)
			return
		}

		err = a.roleUseCase.Delete(v.ID)
		if err != nil {
			a.fastErrorResponses.Execute(w, r, "ROLE_DELETION", http.StatusInternalServerError)
			return
		}

	}

	var response presentation.Response[bool]
	response.Data = true

	encoded, err := json.Marshal(response)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(encoded)
}
