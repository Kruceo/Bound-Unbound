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

	"github.com/gorilla/mux"
)

func Run() {
	fmt.Println("HOST")

	priv, pub := security.GenKeysUseCase{}.GenKeys()
	nodeRepo := adapters.NewInMemoryNodeRepository()
	responseRepo := adapters.NewInMemoryResponseRepository()
	userRepo := adapters.NewFileSystemUserRepo("users.temp.json")

	authMiddleware := middlewares.NewJWTMiddleware("payet").AuthMiddleware

	r := mux.NewRouter()
	apiRouter := routers.SetupNodesRouter(r, &nodeRepo, &responseRepo)
	apiRouter.Use(authMiddleware)

	routers.SetupAuthRouter(r, userRepo, "payet")
	routers.SetupWebsocketRouter(r, &nodeRepo, &responseRepo, priv, pub)

	http.ListenAndServe(":8080", r)
}
