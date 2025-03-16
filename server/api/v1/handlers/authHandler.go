package handlers

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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

	userToken := make([]byte, 64)
	if _, err := rand.Read(userToken); err != nil {
		panic(err)
	}

	userTokenB64 := base64.RawStdEncoding.EncodeToString(userToken)

	type TokenR struct{ Token string }
	response := v1.Response[TokenR]{Data: TokenR{Token: userTokenB64}, Message: ""}

	responseEncoded, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	// expiresAt := time.Now().Add(30 * time.Second)
	aesblock, err := aes.NewCipher([]byte("1234567890123456"))
	if err != nil {
		panic((err))
	}
	c, err := cipher.NewGCM(aesblock)
	if err != nil {
		panic((err))
	}
	encrypted := c.Seal(nil, make([]byte, 12), []byte(time.Now().String()), nil)

	// APIClients[userTokenB64] = APIClient{ExpiresAt: expiresAt}
	w.Header().Add("Set-Cookie", fmt.Sprintf("session=%s;Secure;SameSite", base64.RawStdEncoding.EncodeToString(encrypted)))
	w.Write(responseEncoded)
}

func AuthClientToken(w http.ResponseWriter, r *http.Request) {

	aesblock, _ := aes.NewCipher([]byte("1234567890123456"))
	c, _ := cipher.NewGCM(aesblock)

	if v1.CorsHandler(w, r, "GET, OPTIONS") {
		return
	}
	w.Header().Add("Content-Type", "application/json")
	for _, v := range r.CookiesNamed("session") {
		decoded, _ := base64.RawStdEncoding.DecodeString(v.Value)
		red, _ := c.Open(nil, make([]byte, 12), decoded, nil)
		date, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(red))
		fmt.Println(date.String())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(nil)

	}

	w.Write(nil)
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
