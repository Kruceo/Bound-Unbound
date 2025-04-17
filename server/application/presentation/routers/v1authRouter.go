package routers

import (
	v1 "server2/application/controllers/api/v1"
	"server2/application/infrastructure"
	usecases "server2/application/useCases"

	"github.com/gorilla/mux"
)

func SetupAuthRouter(r *mux.Router, userRepo infrastructure.UserRepository, roleRepo infrastructure.RoleRepository, jwtUseCase *usecases.JwtUseCase) *mux.Router {

	authController := v1.NewV1AuthHandlers(userRepo, roleRepo, jwtUseCase)

	prefix := r.PathPrefix("/auth")

	loginRouter := prefix.Subrouter()

	loginRouter.HandleFunc("/login", authController.AuthLoginHandler).Methods("POST")
	loginRouter.HandleFunc("/token", authController.AuthClientToken).Methods("GET")
	loginRouter.HandleFunc("/register", authController.AuthRegisterHandler).Methods("POST")
	loginRouter.HandleFunc("/register/request", authController.AuthCreateRegisterRequest).Methods("POST")
	loginRouter.HandleFunc("/reset", authController.AuthResetAccountHandler).Methods("POST")
	loginRouter.HandleFunc("/reset/pwd", authController.AuthResetAccountPasswordHandler).Methods("POST")

	return loginRouter
}
