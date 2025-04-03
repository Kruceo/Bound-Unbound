//go:build host
// +build host

package cmd

import (
	"fmt"
	"net/http"
	"server2/application/infrastructure/adapters"
	"server2/application/presentation/middlewares"
	"server2/application/presentation/routers"
	"server2/application/useCases/security"
	"server2/enviroment"

	"github.com/gorilla/mux"
)

func Run() {
	fmt.Println("HOST")

	priv, pub := security.GenKeysUseCase{}.GenKeys()
	nodeRepo := adapters.NewInMemoryNodeRepository()
	responseRepo := adapters.NewInMemoryResponseRepository()
	userRepo := adapters.NewFileSystemUserRepo("users.temp.json")
	authMiddleware := middlewares.NewJWTMiddleware(enviroment.SESSION_SECRET).AuthMiddleware
	corsMiddleware := middlewares.NewCorsMiddleware(enviroment.CORS_ORIGIN, "Authorization", "Content-Type", "Cookie").CorsMiddleware

	r := mux.NewRouter()
	apiRouter := routers.SetupNodesRouter(r, &nodeRepo, &responseRepo)
	authRouter := routers.SetupAuthRouter(r, userRepo, enviroment.SESSION_SECRET)
	routers.SetupWebsocketRouter(r, &nodeRepo, &responseRepo, priv, pub)

	apiRouter.Use(corsMiddleware, authMiddleware)
	authRouter.Use(corsMiddleware)

	http.ListenAndServe(":8080", r)
}
