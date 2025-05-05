package v1

import (
	"server2/application/infrastructure"
	"server2/application/presentation"
	usecases "server2/application/useCases"
)

type V1APIHandlers struct {
	nodePersistenceUseCase *usecases.NodePersistenceUseCase
	responseRepo           infrastructure.ResponsesReporisory
	nodeRoleBindUseCase    *usecases.NodeRoleBindPersistenceUseCase
	fastErrorResponses     presentation.FastErrorResponses
}

func NewV1Handlers(NodeRepo infrastructure.NodeRepository, nodeRoleRepo infrastructure.NodeRoleBindRepository, roleRepo infrastructure.RoleRepository, ResponseRepo infrastructure.ResponsesReporisory) *V1APIHandlers {
	return &V1APIHandlers{nodePersistenceUseCase: usecases.NewNodePersistenceUseCase(NodeRepo), nodeRoleBindUseCase: usecases.NewNodeRoleBindPersistenceUseCase(nodeRoleRepo, usecases.NewNodePersistenceUseCase(NodeRepo), usecases.NewRoleUseCase(roleRepo)), responseRepo: ResponseRepo, fastErrorResponses: presentation.NewFastErrorResponses()}
}
