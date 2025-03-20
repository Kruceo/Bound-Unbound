package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	v1 "unbound-mngr-host/api/v1"
)

type LoginR struct {
	User     string
	Password string
}

type TokenW struct {
	Token string
}

func AuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}
	w.Header().Add("Content-Type", "application/json")
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		v1.FastErrorResponse(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	fmt.Println(string(body))

	var b LoginR
	err = json.Unmarshal(body, &b)

	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	if &b.Password == nil || b.Password == "" || len(b.Password) < 8 {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile("./userdata", os.O_RDONLY, 0600)
	if err != nil {
		v1.FastErrorResponse(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text()
	}

	rawParts := strings.Split(content, ",")
	if len(rawParts) != 2 || rawParts[0] != b.User {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}
	if !v1.VerifyPassword(b.Password, rawParts[1]) {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	jwtToken, err := v1.GenerateJWT(b.User)
	if err != nil {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	t := v1.Response[TokenW]{Data: TokenW{Token: jwtToken}, Message: ""}
	encoded, err := json.Marshal(t)
	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Set-Cookie", "session="+jwtToken+"; SameSite=none; Secure")

	w.Write(encoded)

	fmt.Println("logged")
}

func AuthClientToken(w http.ResponseWriter, r *http.Request) {

	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}

	token, err := v1.JWTMiddleware(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	subject, err := token.Claims.GetSubject()
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Set-Cookie", "user="+subject+"; SameSite=None; Secure")

	w.Write([]byte("Ok"))
}

func AuthRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}

	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		v1.FastErrorResponse(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b LoginR
	err = json.Unmarshal(body, &b)

	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	if &b.Password == nil || len(b.Password) < 8 || len(b.User) < 1 {
		v1.FastErrorResponse(w, r, "BODY_FORMAT", http.StatusBadRequest)
		return
	}

	passwordHash := v1.HashPassword(b.Password)

	file, err := os.OpenFile("./userdata", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		v1.FastErrorResponse(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			v1.FastErrorResponse(w, r, "OVERWRITING_REGISTER", http.StatusUnauthorized)
			return
		}
	}

	file.WriteString(b.User + "," + string(passwordHash))
	w.Write(nil)
}

type HasUserW struct {
	AlreadyRegistered bool
}

func AuthHasUserHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}

	var response v1.Response[HasUserW] = v1.Response[HasUserW]{Message: "", Data: HasUserW{AlreadyRegistered: false}}

	file, err := os.OpenFile("./userdata", os.O_RDONLY, 0600)
	if err == nil {
		defer file.Close()
		response.Data.AlreadyRegistered = true
	}

	responseEncoded, err := json.Marshal(response)
	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	w.Write(responseEncoded)
}
