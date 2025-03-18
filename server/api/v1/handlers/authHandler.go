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

func AuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}
	w.Header().Add("Content-Type", "application/json")
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("body read error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	var b LoginR
	err = json.Unmarshal(body, &b)

	if err != nil {
		fmt.Println("json decode error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	if &b.Password == nil || b.Password == "" || len(b.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}

	file, err := os.OpenFile("./userdata", os.O_RDONLY, 0600)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text()
	}

	rawParts := strings.Split(content, ",")
	if len(rawParts) != 2 && rawParts[0] != b.User {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(nil)
		return
	}
	if !v1.VerifyPassword(b.Password, rawParts[1]) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(nil)
		return
	}

	jwtToken, err := v1.GenerateJWT(b.User)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	t := v1.Response[string]{Data: jwtToken, Message: ""}
	encoded, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	w.Write(encoded)
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

	subject, _ := token.Claims.GetSubject()

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
		fmt.Println("body read error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	var b LoginR
	err = json.Unmarshal(body, &b)

	if err != nil {
		fmt.Println("json decode error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	if &b.Password == nil || b.Password == "" || len(b.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}

	passwordHash := v1.HashPassword(b.Password)

	file, err := os.OpenFile("./userdata", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(nil)
			return
		}
	}

	file.WriteString(b.User + "," + string(passwordHash))
	w.Write(nil)
}
