package routers

import (
	"crypto/ecdh"
	"server2/application/controllers"
	"server2/application/infrastructure"

	"github.com/gorilla/mux"
)

func SetupWebsocketRouter(r *mux.Router, nodeRepo infrastructure.NodeRepository, responseRepo infrastructure.ResponsesReporisory, privateKey *ecdh.PrivateKey, publicKey *ecdh.PublicKey) {
	websocketController := controllers.NewHostController(nodeRepo, responseRepo, *privateKey, *publicKey)
	router := r.PathPrefix("/ws").Subrouter()
	router.HandleFunc("/node", websocketController.OnMessageHandler)
}
