//go:build host
// +build host

package runners

import (
	"fmt"
	"net/http"
	"server2/application/adapters"
	"server2/application/controllers"
	"server2/application/controllers/api/v1/endpoints"
	"server2/application/useCases/security"
)

func Run() {
	fmt.Println("HOST")
	priv, pub := security.GenKeysUseCase{}.GenKeys()
	nodeRepo := adapters.NewInMemoryNodeRepository()
	responseRepo := adapters.NewInMemoryResponseRepository()
	websocketController := controllers.NewHostController(&nodeRepo, &responseRepo, *priv, *pub)
	v1Handlers := endpoints.V1Handlers{NodeRepo: &nodeRepo, ResponseRepo: &responseRepo}

	http.HandleFunc("/ws", websocketController.OnMessageHandler)

	http.HandleFunc("/auth/login", endpoints.AuthLoginHandler)

	http.HandleFunc("/auth/token", endpoints.AuthClientToken)

	http.HandleFunc("/auth/register", endpoints.AuthRegisterHandler)

	http.HandleFunc("/auth/status", endpoints.AuthHasUserHandler)

	http.HandleFunc("/auth/reset", endpoints.AuthResetAccountHandler)

	http.HandleFunc("/v1/connections", v1Handlers.ConnectionsHandler)

	http.HandleFunc("/v1/connections/{connection}/blocks", v1Handlers.BlockAddressHandler)

	http.HandleFunc("/v1/connections/{connection}/redirects", v1Handlers.RedirectAddressHandler)

	http.HandleFunc("/v1/connections/{connection}/reload", v1Handlers.ReloadHandler)

	http.HandleFunc("/v1/connections/{connection}/confighash", v1Handlers.ConfigHashHandler)
	http.ListenAndServe(":8080", nil)
}
