package routers

import (
	v1 "server2/application/controllers/api/v1"
	"server2/application/infrastructure"

	"github.com/gorilla/mux"
)

func SetupAuthRouter(r *mux.Router, userRepo infrastructure.UserRepository, roleRepo infrastructure.RoleRepository, sessionSecret string) *mux.Router {

	authController := v1.NewV1AuthHandlers(userRepo, roleRepo, sessionSecret)

	router := r.PathPrefix("/auth").Subrouter()
	router.HandleFunc("/login", authController.AuthLoginHandler).Methods("POST")
	router.HandleFunc("/token", authController.AuthClientToken).Methods("GET")
	router.HandleFunc("/register", authController.AuthRegisterHandler).Methods("POST")
	router.HandleFunc("/register/request", authController.AuthCreateRegisterRequest).Methods("POST")

	router.HandleFunc("/roles", authController.AuthGetRoles).Methods("GET")
	router.HandleFunc("/roles", authController.AuthPostRole).Methods("POST")
	router.HandleFunc("/roles", authController.AuthDeleteRole).Methods("DELETE")

	router.HandleFunc("/reset", authController.AuthResetAccountHandler).Methods("POST")
	router.HandleFunc("/reset/pwd", authController.AuthResetAccountPasswordHandler).Methods("POST")

	return router
}
