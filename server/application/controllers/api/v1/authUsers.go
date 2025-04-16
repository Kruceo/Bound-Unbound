package v1

import (
	"encoding/json"
	"net/http"
	"server2/application/presentation"
)

type InnerRole struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type GetUsersW struct {
	Name string    `json:"name"`
	Role InnerRole `json:"role"`
}

func (a *v1AuthHandlers) AuthGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.userUseCase.SearchByName(".*")
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "USERS_RECOVERY", http.StatusInternalServerError)
		return
	}

	var response presentation.Response[[]GetUsersW]

	for _, v := range users {
		userRole, err := a.roleUseCase.Get(v.RoleID)
		if err != nil {
			a.fastErrorResponses.Execute(w, r, "ROLE_RECOVERY", http.StatusInternalServerError)
			return
		}

		response.Data = append(response.Data, GetUsersW{
			Name: v.Username,
			Role: InnerRole{
				ID:          userRole.ID,
				Name:        userRole.Name,
				Permissions: userRole.Permissions,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
