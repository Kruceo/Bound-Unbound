package v1

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	w.Header().Add("Content-Type", "application/json")
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
		fmt.Println(err)
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

	w.Header().Add("Set-Cookie", "session="+jwtToken+"; SameSite=none; Secure")

	w.Write(encoded)
}

func (a *v1AuthHandlers) AuthClientToken(w http.ResponseWriter, r *http.Request) {
	_, err := a.jwtManager.JWTMiddleware(r)
	if err != nil {
		fmt.Println(err)
		a.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.Write([]byte("Ok"))
}

type RegisterW struct {
	SecretCode string `json:"secretCode"`
}

func (a *v1AuthHandlers) AuthRegisterHandler(w http.ResponseWriter, r *http.Request) {

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

	if &b.Password == nil || len(b.Password) < 8 || len(b.User) < 1 {
		a.fastErrorResponses.Execute(w, r, "BODY_FORMAT", http.StatusBadRequest)
		return
	}

	passwordHash := a.hashPassword.Hash(b.Password)

	extraSecretCode := make([]byte, 12)
	rand.Read(extraSecretCode)

	extraSecretCodeB64 := base64.RawStdEncoding.EncodeToString(extraSecretCode) // real string (will show it to client)
	hashedSecretCode := a.hashPassword.Hash(extraSecretCodeB64)                 // hashed (store it)

	_, err = a.userRepo.Save(b.User, string(passwordHash), 0, string(hashedSecretCode))

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

type HasUserW struct {
	AlreadyRegistered bool `json:"alreadyRegistered"`
}

func (a *v1AuthHandlers) AuthHasUserHandler(w http.ResponseWriter, r *http.Request) {

	var response presentation.Response[HasUserW] = presentation.Response[HasUserW]{Message: "", Data: HasUserW{AlreadyRegistered: false}}

	user, err := a.userRepo.FindOneByName(".*")

	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	if user != nil {
		response.Data.AlreadyRegistered = true
	}

	responseEncoded, err := json.Marshal(response)
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	w.Write(responseEncoded)
}

type ResetAccountR struct {
	User       string `json:"user"`
	SecretCode string `json:"secretCode"`
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

	responseEncoded, err := json.Marshal(presentation.Response[bool]{Message: "ok", Data: true})
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	// do anithing to change the user
	// ...
	// ...
	panic("not implemented")
	if err != nil {
		a.fastErrorResponses.Execute(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	w.Write(responseEncoded)
}
