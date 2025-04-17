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
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Role InnerRole `json:"role"`
}

func (a *V1AdminHandlers) AdminGetUsers(w http.ResponseWriter, r *http.Request) {
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
			ID:   v.ID,
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
}

type DeleteUserR struct {
	ID string `json:"id"`
}

func (a *V1AdminHandlers) AdminDeleteUsers(w http.ResponseWriter, r *http.Request) {
	requesterUser, err := a.getUserFromJWTBearerUseCase.Execute(r.Header.Get("authorization"))
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "USER_RECOVERY", http.StatusInternalServerError)
		return
	}

	var req []DeleteUserR
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	for _, v := range req {
		if v.ID == requesterUser.ID {
			a.fastErrorResponses.Execute(w, r, "INVALID_USER", http.StatusBadRequest)
			return
		}
	}

	for _, v := range req {
		err = a.userUseCase.Delete(v.ID)
		if err != nil {
			a.fastErrorResponses.Execute(w, r, "USER_RENOVE", http.StatusInternalServerError)
			return
		}
	}
}
