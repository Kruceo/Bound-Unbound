//go:build host
// +build host

package cmd

import (
	"fmt"
	"net/http"
	"server2/application/infrastructure/adapters"
	"server2/application/presentation/middlewares"
	"server2/application/presentation/routers"
	usecases "server2/application/useCases"
	"server2/application/useCases/security"
	"server2/enviroment"
	"server2/utils"

	"github.com/gorilla/mux"
)

func Run() {
	priv, pub := security.GenKeysUseCase{}.GenKeys()
	nodeRepo := adapters.NewInMemoryNodeRepository()
	responseRepo := adapters.NewInMemoryResponseRepository()
	userRepo := adapters.NewFileSystemUserRepo("users.temp.json")
	roleRepo := adapters.NewFileRoleRepository("roles.temp.json")

	jwtUseCase := usecases.NewJWTUseCase(enviroment.SESSION_SECRET)

	authMiddleware := middlewares.NewJWTMiddleware(enviroment.SESSION_SECRET).AuthMiddleware
	corsMiddleware := middlewares.NewCorsMiddleware(enviroment.CORS_ORIGIN, "Authorization", "Content-Type", "Cookie").CorsMiddleware
	bforceMiddleware := middlewares.NewBruteForceMiddleware().BruteForceMiddleware
	roleMiddleware := middlewares.NewRoleMiddleware([]string{"admin"}, jwtUseCase, userRepo, roleRepo).Middleware

	r := mux.NewRouter()
	apiRouter := routers.SetupNodesRouter(r, &nodeRepo, &responseRepo)
	loginRouter := routers.SetupAuthRouter(r, userRepo, roleRepo, jwtUseCase)
	adminRouter := routers.SetupAdminRouter(r, userRepo, roleRepo, jwtUseCase)
	routers.SetupWebsocketRouter(r, &nodeRepo, &responseRepo, priv, pub)

	apiRouter.Use(corsMiddleware, authMiddleware)
	loginRouter.Use(corsMiddleware, bforceMiddleware)
	adminRouter.Use(corsMiddleware, authMiddleware, roleMiddleware)
	port := utils.GetEnvOrDefaultNumber("port", 8080)
	fmt.Println("listening", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
