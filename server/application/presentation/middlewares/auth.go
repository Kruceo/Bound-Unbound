package middlewares

import (
	"fmt"
	"net/http"
	"server2/application/presentation"
	usecases "server2/application/useCases"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMiddleware struct {
	jwtManager         usecases.JwtUseCase
	fastErrorResponses presentation.FastErrorResponses
}

func NewJWTMiddleware(sessionSecret string) *JWTMiddleware {
	return &JWTMiddleware{jwtManager: *usecases.NewJWTUseCase(sessionSecret), fastErrorResponses: presentation.NewFastErrorResponses()}
}

// verify if authorization is a valid jwt, and if jwt address is compatible with requester ip
func (j *JWTMiddleware) test(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {

	authorization, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	token, err := j.jwtManager.ValidateJWT(string(authorization))

	if err != nil {
		return nil, err
	}

	subject, err := token.Claims.GetSubject()

	if err != nil {
		return nil, err
	}

	user, userAddress := j.jwtManager.ParseJWTSubject(subject)

	if userAddress != strings.Split(r.RemoteAddr, ":")[0] {
		j.fastErrorResponses.Execute(w, r, "INVALID_ADDRESS", http.StatusUnauthorized)
		return nil, fmt.Errorf("jwt address does not match with request address")
	}

	r.Header.Add("user-id", user)
	// r.Header.Add("role-id", user)

	return token, nil
}

func (j *JWTMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens, err := j.test(w, r)
		if err != nil {
			j.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
			return
		}
		subject, err := tokens.Claims.GetSubject()
		if err != nil {
			j.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
			return
		}

		userid, originalAddress := j.jwtManager.ParseJWTSubject(subject)
		r.Header.Set("user-id", userid)
		r.Header.Set("user-original-address", originalAddress)
		next.ServeHTTP(w, r)
	})
}
