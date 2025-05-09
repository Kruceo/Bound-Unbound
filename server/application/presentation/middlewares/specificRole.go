package middlewares

import (
	"fmt"
	"net/http"
	"server2/application/infrastructure"
	"server2/application/presentation"
	usecases "server2/application/useCases"
	"slices"
	"strings"
)

type RoleMiddleware struct {
	getUserFromBearerUseCase *usecases.GetUserFromJWTBearerUseCase
	rolesUseCase             *usecases.RoleUseCase
	fastErrorResponses       *presentation.FastErrorResponses
	permissions              []string
}

func NewRoleMiddleware(permissions []string, jwt *usecases.JwtUseCase, userRepo infrastructure.UserRepository, roleRepo infrastructure.RoleRepository) *RoleMiddleware {
	f := presentation.NewFastErrorResponses()
	return &RoleMiddleware{
		permissions:              permissions,
		rolesUseCase:             usecases.NewRoleUseCase(roleRepo),
		getUserFromBearerUseCase: usecases.NewGetUserFromJWTBearerUseCase(userRepo, jwt),
		fastErrorResponses:       &f,
	}
}
func (rm *RoleMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requesterUser, err := rm.getUserFromBearerUseCase.Execute(r.Header.Get("Authorization"))
		if err != nil {
			rm.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
			return
		}

		userRole, err := rm.rolesUseCase.Get(requesterUser.RoleID)
		if err != nil {
			rm.fastErrorResponses.Execute(w, r, "RECOVER_PERMISSION", http.StatusUnauthorized)
			return
		}

		if slices.Contains(rm.permissions, "*") {
			fmt.Println("any perm can reach this route")
			r.Header.Set("X-Role-ID", requesterUser.RoleID)
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("just perm", strings.Join(rm.permissions, ","))
		permSet := make(map[string]struct{})
		for _, perm := range userRole.Permissions {
			permSet[perm] = struct{}{}
		}

		for _, requiredPerm := range rm.permissions {
			if _, ok := permSet[requiredPerm]; !ok {
				rm.fastErrorResponses.Execute(w, r, "AUTH", http.StatusUnauthorized)
				return
			}
		}

		r.Header.Set("X-Role-ID", requesterUser.RoleID)
		next.ServeHTTP(w, r)
	})
}
