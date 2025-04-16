package middlewares

import (
	"net/http"
	"server2/application/infrastructure"
	"server2/application/presentation"
	usecases "server2/application/useCases"
)

type RoleMiddleware struct {
	getUserFromBearerUseCase *usecases.GetUserFromJWTBearerUseCase
	fastErrorResponses       *presentation.FastErrorResponses
	permissions              []string
}

func NewRoleMiddleware(permissions []string, jwt *usecases.JwtUseCase, userRepo infrastructure.UserRepository, roleRepo infrastructure.RoleRepository) *RoleMiddleware {
	return &RoleMiddleware{
		permissions:              permissions,
		getUserFromBearerUseCase: usecases.NewGetUserFromJWTBearerUseCase(userRepo, jwt),
	}
}

func (rm *RoleMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			requesterUser, err := rm.getUserFromBearerUseCase.Execute(r.Header.Get("Authorization"))
			if err != nil {
				rm.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
				return
			}

			if !requesterUser.IsAdmin() {
				rm.fastErrorResponses.Execute(w, r, "NO_PERMISSION", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
}
