package routers

import (
	v1 "server2/application/controllers/api/v1"
	"server2/application/infrastructure"

	"github.com/gorilla/mux"
)

func SetupAuthRouter(r *mux.Router, userRepo infrastructure.UserRepository, sessionSecret string) *mux.Router {

	authController := v1.NewV1AuthHandlers(userRepo, sessionSecret)

	router := r.PathPrefix("/auth").Subrouter()
	router.HandleFunc("/login", authController.AuthLoginHandler).Methods("POST")
	router.HandleFunc("/token", authController.AuthClientToken).Methods("GET")
	router.HandleFunc("/register", authController.AuthRegisterHandler).Methods("POST")
	router.HandleFunc("/status", authController.AuthHasUserHandler).Methods("GET")
	router.HandleFunc("/reset", authController.AuthResetAccountHandler).Methods("POST")
	return router
}
