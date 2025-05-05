package v1

import (
	"server2/application/infrastructure"
	"server2/application/presentation"
	usecases "server2/application/useCases"
)

type V1AdminHandlers struct {
	roleUseCase                 *usecases.RoleUseCase
	userUseCase                 *usecases.UserUseCase
	getUserFromJWTBearerUseCase *usecases.GetUserFromJWTBearerUseCase
	bindsPersistence            *usecases.NodeRoleBindPersistenceUseCase
	hashPassword                *usecases.PassowrdHashUseCase
	jwtManager                  *usecases.JwtUseCase
	fastErrorResponses          presentation.FastErrorResponses
}

func NewV1AdminHandlers(userRepo infrastructure.UserRepository, roleRepo infrastructure.RoleRepository, nodeRepo infrastructure.NodeRepository, nrbRepo infrastructure.NodeRoleBindRepository, jwtUseCase *usecases.JwtUseCase) *V1AdminHandlers {
	pwMan := usecases.NewPassowrdHashUseCase()
	return &V1AdminHandlers{
		roleUseCase:                 usecases.NewRoleUseCase(roleRepo),
		userUseCase:                 usecases.NewUserUseCase(userRepo),
		bindsPersistence:            usecases.NewNodeRoleBindPersistenceUseCase(nrbRepo, nodeRepo, roleRepo),
		getUserFromJWTBearerUseCase: usecases.NewGetUserFromJWTBearerUseCase(userRepo, jwtUseCase),
		jwtManager:                  jwtUseCase,
		hashPassword:                &pwMan,
		fastErrorResponses:          presentation.NewFastErrorResponses(),
	}
}
