package routers

import (
	v1 "server2/application/controllers/api/v1"
	"server2/application/infrastructure"
	usecases "server2/application/useCases"

	"github.com/gorilla/mux"
)

func SetupAdminRouter(r *mux.Router, userRepo infrastructure.UserRepository, roleRepo infrastructure.RoleRepository, nodeRepo infrastructure.NodeRepository, nrbRepo infrastructure.NodeRoleBindRepository, jwtUseCase *usecases.JwtUseCase) *mux.Router {

	adminController := v1.NewV1AdminHandlers(userRepo, roleRepo, nodeRepo, nrbRepo, jwtUseCase)

	adminRouter := r.PathPrefix("/admin").Subrouter()

	adminRouter.HandleFunc("/roles", adminController.AdminGetRoles).Methods("GET")
	adminRouter.HandleFunc("/roles", adminController.AdminPostRole).Methods("POST")
	adminRouter.HandleFunc("/roles", adminController.AdminDeleteRole).Methods("DELETE")

	adminRouter.HandleFunc("/roles/bind/nodes", adminController.AdminPostBinds).Methods("POST")
	adminRouter.HandleFunc("/roles/bind/nodes", adminController.AdminGetBinds).Methods("GET")
	adminRouter.HandleFunc("/roles/bind/nodes", adminController.AdminDeleteBinds).Methods("DELETE")

	adminRouter.HandleFunc("/users", adminController.AdminGetUsers).Methods("GET")
	adminRouter.HandleFunc("/users", adminController.AdminDeleteUsers).Methods("DELETE")
	return adminRouter
}
