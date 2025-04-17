package v1

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"server2/application/presentation"
)

type LoginR struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type TokenW struct {
	Token string `json:"token"`
}

func (a *v1AuthHandlers) AuthLoginHandler(w http.ResponseWriter, r *http.Request) {

	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b LoginR
	err = json.Unmarshal(body, &b)

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	// verify if received password is ok
	if &b.Password == nil || b.Password == "" || len(b.Password) < 8 {
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusBadRequest)

		return
	}

	storedUser, err := a.userRepo.FindOneByName(b.User)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}
	if !a.hashPassword.VerifyPassword(b.Password, storedUser.Password) {
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	jwtToken, err := a.jwtManager.GenerateJWT(b.User, r.RemoteAddr)
	if err != nil {
		fmt.Println(err)
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	t := presentation.Response[TokenW]{Data: TokenW{Token: jwtToken}, Message: ""}
	encoded, err := json.Marshal(t)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", "session="+jwtToken+"; SameSite=none; Secure")
	w.Write(encoded)
}

func (a *v1AuthHandlers) AuthClientToken(w http.ResponseWriter, r *http.Request) {
	_, err := a.jwtManager.TokenFromBearer(r.Header.Get("Authorization"))

	if err != nil {
		fmt.Println(err)
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	// user, err := a.userRepo.FindOneByName(".*")

	// if err != nil || user == nil {
	// 	a.fastErrorResponses.Execute(w, r, "NO_USERS", http.StatusUnauthorized)
	// 	return
	// }

	w.Header().Add("Content-Type", "application/json")

	w.Write([]byte("Ok"))
}

type RegisterR struct {
	User     string `json:"user"`
	Password string `json:"password"`
	RouteId  string `json:"routeId"`
}

type RegisterW struct {
	SecretCode string `json:"secretCode"`
}

func (a *v1AuthHandlers) AuthRegisterHandler(w http.ResponseWriter, r *http.Request) {
	userCount, err := a.userRepo.Count()
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "ALREADY_REGISTERED", http.StatusInternalServerError)
		return
	}

	var body []byte
	body, err = io.ReadAll(r.Body)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b RegisterR
	err = json.Unmarshal(body, &b)

	var newUserRoleID string = "0"
	if userCount > 0 {
		// change this to use a routeRepe entry
		// first registered user is a admin
		// and this admin can create route repo entries to other users register in this

		storedRoleID, routeExists := a.routesRepo.Exists(b.RouteId)
		if !routeExists {
			a.fastErrorResponses.Execute(w, r, "NOT_FOUND", http.StatusNotFound)
			return
		}
		if err != nil {
			// handle the error
			a.fastErrorResponses.Execute(w, r, "ROLE_PARSE_ERROR", http.StatusInternalServerError)
			return
		}

		newUserRoleID = storedRoleID
	}

	if _, err = a.roleUseCase.Get(newUserRoleID); err != nil {
		a.fastErrorResponses.Execute(w, r, "ROLE_NOT_FOUND", http.StatusBadRequest)
		return
	}

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	if &b.Password == nil || len(b.Password) < 8 || len(b.User) < 1 {
		a.fastErrorResponses.Execute(w, r, "BODY_FORMAT", http.StatusBadRequest)
		return
	}

	passwordHash := a.hashPassword.Hash(b.Password)

	extraSecretCode := make([]byte, 12)
	rand.Read(extraSecretCode)

	extraSecretCodeB64 := base64.RawStdEncoding.EncodeToString(extraSecretCode) // real string (will show it to client)
	hashedSecretCode := a.hashPassword.Hash(extraSecretCodeB64)                 // hashed (store it)

	_, err = a.userRepo.Save(b.User, string(passwordHash), newUserRoleID, string(hashedSecretCode))

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "REPO", http.StatusInternalServerError)
		return
	}

	var res presentation.Response[RegisterW] = presentation.Response[RegisterW]{Message: "", Data: RegisterW{SecretCode: extraSecretCodeB64}}
	encodedRes, err := json.Marshal(res)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Write(encodedRes)
}

type CreateRegisterR struct {
	RoleID string `json:"roleId"`
}

type CreateRegisterW struct {
	RouteID string `json:"routeId"`
}

func (a *v1AuthHandlers) AuthCreateRegisterRequest(w http.ResponseWriter, r *http.Request) {
	requesterUser, err := a.getUserFromJWTBearerUseCase.Execute(r.Header.Get("Authorization"))
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	if !requesterUser.IsAdmin() {
		a.fastErrorResponses.Execute(w, r, "NOT_ADMIN", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b CreateRegisterR
	err = json.Unmarshal(body, &b)

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	routeID, err := a.routesRepo.Gen(b.RoleID)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "ROUTE_GEN", http.StatusInternalServerError)
		return
	}

	var response presentation.Response[CreateRegisterW] = presentation.Response[CreateRegisterW]{Data: CreateRegisterW{RouteID: routeID}}

	responseEncoded, err := json.Marshal(response)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseEncoded)
}

type ResetAccountR struct {
	User       string `json:"user"`
	SecretCode string `json:"secretCode"`
}

type ResetAccountW struct {
	RouteId string `json:"routeId"`
}

func (a *v1AuthHandlers) AuthResetAccountHandler(w http.ResponseWriter, r *http.Request) {

	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b ResetAccountR
	err = json.Unmarshal(body, &b)

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	user, err := a.userRepo.FindOneByName(b.User)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	if !a.hashPassword.VerifyPassword(b.SecretCode, user.RecoveryCode) {
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	routeId, err := a.routesRepo.Gen(user.ID + "," + r.RemoteAddr)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "GENROUTE", http.StatusInternalServerError)
		return
	}

	responseEncoded, err := json.Marshal(presentation.Response[ResetAccountW]{Message: "ok", Data: ResetAccountW{RouteId: routeId}})
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	w.Write(responseEncoded)
}

type NewPasswordR struct {
	RouteId  string `json:"routeId"`
	Password string `json:"password"`
}

func (a *v1AuthHandlers) AuthResetAccountPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b NewPasswordR
	err = json.Unmarshal(body, &b)

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	info, routeExists := a.routesRepo.Exists(b.RouteId)
	if !routeExists {
		a.fastErrorResponses.Execute(w, r, "NOT_FOUND", http.StatusNotFound)
		return
	}

	splited := strings.Split(info, ",")
	userId, originalAddr := splited[0], splited[1]

	if originalAddr != r.RemoteAddr {
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}
	user, err := a.userRepo.Get(userId)
	if err != nil {
		fmt.Println(err)
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	passwordHash := a.hashPassword.Hash(b.Password)

	extraSecretCode := make([]byte, 12)
	rand.Read(extraSecretCode)

	extraSecretCodeB64 := base64.RawStdEncoding.EncodeToString(extraSecretCode) // real string (will show it to client)
	hashedSecretCode := a.hashPassword.Hash(extraSecretCodeB64)                 // hashed (store it)

	err = a.userRepo.Update(userId, user.Username, string(passwordHash), user.RoleID, string(hashedSecretCode))
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "UPDATE_USER", http.StatusInternalServerError)
		return
	}

	responseEncoded, err := json.Marshal(presentation.Response[RegisterW]{Message: "ok", Data: RegisterW{SecretCode: extraSecretCodeB64}})
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	w.Write(responseEncoded)
}
