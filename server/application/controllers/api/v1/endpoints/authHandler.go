package endpoints

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	v1 "server2/application/controllers/api/v1"

	"strings"
	"time"
)

type LoginR struct {
	User     string
	Password string
}

type TokenW struct {
	Token string
}

type BlockedRequester struct {
	LastTry   time.Time
	LimitTime time.Time
}

var storedRequests = make(map[string]*BlockedRequester)

func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			for i, v := range storedRequests {
				if time.Now().UnixMilli()-v.LastTry.UnixMilli() > 1000*60 {
					// fmt.Println("deleting", i)
					delete(storedRequests, i)
				}
			}
			// fmt.Println("Blocked Requesters cleaning")
		}
	}()
}

func AuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}

	now := time.Now()

	if timestamp, exists := storedRequests[r.RemoteAddr]; exists {
		now := time.Now()
		if timestamp.LimitTime.UnixMilli() > now.UnixMilli() {
			timestamp.LimitTime = timestamp.LimitTime.Add(15 * time.Second)
			timestamp.LastTry = now
			v1.FastErrorResponse(w, r, "AUTH_BLOCKED", http.StatusUnauthorized)
			fmt.Println("blocked", r.RemoteAddr, timestamp.LimitTime.String())
			return
		}
	}

	newLimitTime := BlockedRequester{LastTry: now, LimitTime: now.Add(2 * time.Second)}
	storedRequests[r.RemoteAddr] = &newLimitTime

	w.Header().Add("Content-Type", "application/json")
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

	if &b.Password == nil || b.Password == "" || len(b.Password) < 8 {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile("./userdata", os.O_RDONLY, 0600)
	if err != nil {
		v1.FastErrorResponse(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		v1.FastErrorResponse(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	rawParts := strings.Split(string(content), ",")
	if len(rawParts) != 3 || rawParts[0] != b.User {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}
	if !v1.VerifyPassword(b.Password, rawParts[1]) {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	jwtToken, err := v1.GenerateJWT(b.User, r.RemoteAddr)
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

	_, err := v1.JWTMiddleware(w, r)
	if err != nil {
		fmt.Println(err)
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.Write([]byte("Ok"))
}

type RegisterW struct {
	SecretCode string
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

	extraSecretCode := make([]byte, 12)
	rand.Read(extraSecretCode)

	extraSecretCodeB64 := base64.RawStdEncoding.EncodeToString(extraSecretCode) // real string (will show it to client)
	hashedSecretCode := v1.HashPassword(extraSecretCodeB64)                     // hashed (store it)
	file.WriteString(fmt.Sprintf("%s,%s,%s", b.User, string(passwordHash), hashedSecretCode))

	var res v1.Response[RegisterW] = v1.Response[RegisterW]{Message: "", Data: RegisterW{SecretCode: extraSecretCodeB64}}
	encodedRes, err := json.Marshal(res)
	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_ENCODING", http.StatusInternalServerError)
		return
	}
	w.Write(encodedRes)
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

func AuthResetAccountHandler(w http.ResponseWriter, r *http.Request) {
	if v1.CorsHandler(w, r, "POST, OPTIONS") {
		return
	}

	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		v1.FastErrorResponse(w, r, "READ_BODY", http.StatusInternalServerError)
		return
	}

	var b RegisterW
	err = json.Unmarshal(body, &b)

	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_DECODE", http.StatusInternalServerError)
		return
	}

	file, err := os.OpenFile("./userdata", os.O_RDONLY, 0600)
	if err != nil {
		v1.FastErrorResponse(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		v1.FastErrorResponse(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	rawParts := strings.Split(string(content), ",")
	if len(rawParts) != 3 {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}
	if !v1.VerifyPassword(b.SecretCode, rawParts[2]) {
		v1.FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return
	}

	responseEncoded, err := json.Marshal(v1.Response[bool]{Message: "ok", Data: true})
	if err != nil {
		v1.FastErrorResponse(w, r, "JSON_ENCODE", http.StatusInternalServerError)
		return
	}

	err = os.Remove("userdata")
	if err != nil {
		v1.FastErrorResponse(w, r, "LOGIN", http.StatusInternalServerError)
		return
	}

	w.Write(responseEncoded)
}
