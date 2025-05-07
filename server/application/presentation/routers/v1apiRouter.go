package routers

import (
	v1 "server2/application/controllers/api/v1"
	"server2/application/infrastructure"

	"github.com/gorilla/mux"
)

func SetupNodesRouter(r *mux.Router, nodeRepo infrastructure.NodeRepository, responseRepo infrastructure.ResponsesReporisory, nodeRoleRepo infrastructure.NodeRoleBindRepository, roleRepo infrastructure.RoleRepository) *mux.Router {

	apiController := v1.NewV1Handlers(nodeRepo, nodeRoleRepo, roleRepo, responseRepo)
	router := r.PathPrefix("/v1").Subrouter()

	router.HandleFunc("/connections", apiController.ConnectionsHandler).Methods("GET")
	router.HandleFunc("/connections/{connection}/blocks", apiController.BlockAddressHandler)
	router.HandleFunc("/connections/{connection}/redirects", apiController.RedirectAddressHandler)
	router.HandleFunc("/connections/{connection}/reload", apiController.ReloadHandler)
	router.HandleFunc("/connections/{connection}/confighash", apiController.ConfigHashHandler)

	return router
}
