package v1

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

const (
	times           = 3
	availableMemomy = 64 * 1024
	threads         = 2
	saltLen         = 32
)

func parseParams(timeStr, memoryStr, threadsStr, keyLenStr string) (uint32, uint32, uint8, uint32) {
	var time, memory, keyLen uint32
	var threads uint8
	fmt.Sscanf(timeStr, "%d", &time)
	fmt.Sscanf(memoryStr, "%d", &memory)
	fmt.Sscanf(threadsStr, "%d", &threads)
	fmt.Sscanf(keyLenStr, "%d", &keyLen)
	return time, memory, threads, keyLen
}

// Output is base64
func HashPassword(password string) []byte {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	hash := argon2.IDKey([]byte(password), salt, times, availableMemomy, threads, saltLen)
	// fmt.Println(salt)
	hashB64, saltB64 := []byte(base64.RawStdEncoding.EncodeToString(hash)), []byte(base64.RawStdEncoding.EncodeToString(salt))

	// fmt.Printf("hashb64 %s\n", hashB64)
	// fmt.Printf("saltb64 %s\n", saltB64)
	hashString := fmt.Sprintf("%d$%d$%d$%d$%s$%s",
		times, availableMemomy, threads, saltLen,
		saltB64,
		hashB64,
	)

	return []byte(hashString)
}

// input is base64
func VerifyPassword(password string, hash string) bool {

	// time$availableMemomy$threads$saltLen$salt(b64)$hash(b64)
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}
	time, memory, threads, keyLen := parseParams(parts[0], parts[1], parts[2], parts[3])
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		panic(err)
	}
	hashStored, _ := base64.RawStdEncoding.DecodeString(parts[5])

	// println("time", time)
	// println("mem", memory)
	// println("threads", threads)
	// println("keylen", keyLen)
	// fmt.Printf("salt %x %d\n", salt, len(salt))
	// fmt.Printf("hash %x\n", hashStored)

	hashComputed := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	// fmt.Printf("hashcomputed %x\n", hashComputed)

	return subtle.ConstantTimeCompare(hashStored, hashComputed) == 1
}

var sessionSecret []byte

func InitAuth() {
	sessionSecret = []byte(os.Getenv("SESSION_SECRET"))
	if len(sessionSecret) < 10 {
		sessionSecret = make([]byte, 24)
		fmt.Println("SESSION_SECRET have less than 10 bytes, replacing with random bytes.")
		_, err := rand.Read(sessionSecret)
		if err != nil {
			panic(err)
		}
	}
}

func GenerateJWT(username string, address string) (string, error) {
	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%s$%s", username, address),
		"exp": time.Now().Add(time.Second * 3600).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(sessionSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWTSubject(subject string) (string, string) {
	splited := strings.Split(subject, "$")
	return splited[0], strings.Split(splited[1], ":")[0]
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
		}
		return sessionSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func JWTMiddleware(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {

	w.Header().Set("Content-Type", "application/json")

	authorization, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	token, err := ValidateJWT(string(authorization))

	if err != nil {
		FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)

		return nil, err
	}

	subject, err := token.Claims.GetSubject()

	if err != nil {
		fmt.Println(err)
		FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return nil, err
	}

	user, userAddress := ParseJWTSubject(subject)

	if userAddress != strings.Split(r.RemoteAddr, ":")[0] {
		FastErrorResponse(w, r, "AUTH", http.StatusUnauthorized)
		return nil, fmt.Errorf("jwt address does not match with request address")
	}
	w.Header().Add("Set-Cookie", "user="+user+"; SameSite=Lax;")

	return token, nil
}
