package v1

import (
	"server2/application/infrastructure"
	"server2/application/infrastructure/adapters"
	"server2/application/presentation"
	usecases "server2/application/useCases"
)

type v1AuthHandlers struct {
	roleUseCase                 *usecases.RoleUseCase
	userUseCase                 *usecases.UserUseCase
	getUserFromJWTBearerUseCase *usecases.GetUserFromJWTBearerUseCase
	userRepo                    infrastructure.UserRepository
	routesRepo                  infrastructure.RoutesRepository
	hashPassword                usecases.PassowrdHashUseCase
	jwtManager                  *usecases.JwtUseCase
	fastErrorResponses          presentation.FastErrorResponses
}

func NewV1AuthHandlers(userRepo infrastructure.UserRepository, roleRepo infrastructure.RoleRepository, sessionSecret string) *v1AuthHandlers {
	jwtUseCase := usecases.NewJWTUseCase(sessionSecret)
	return &v1AuthHandlers{
		roleUseCase:                 usecases.NewRoleUseCase(roleRepo),
		userUseCase:                 usecases.NewUserUseCase(userRepo),
		getUserFromJWTBearerUseCase: usecases.NewGetUserFromJWTBearerUseCase(userRepo, jwtUseCase),
		userRepo:                    userRepo,
		jwtManager:                  jwtUseCase,
		hashPassword:                usecases.NewPassowrdHashUseCase(),
		fastErrorResponses:          presentation.NewFastErrorResponses(),
		routesRepo:                  adapters.NewInMemoryRoutesRepository(),
	}
}
