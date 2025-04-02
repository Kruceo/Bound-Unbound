package v1

import (
	"server2/application/infrastructure"
	"server2/application/presentation"
	usecases "server2/application/useCases"
)

type v1AuthHandlers struct {
	userRepo           infrastructure.UserRepository
	hashPassword       usecases.PassowrdHashUseCase
	jwtManager         *usecases.JwtUseCase
	fastErrorResponses presentation.FastErrorResponses
}

func NewV1AuthHandlers(userRepo infrastructure.UserRepository, sessionSecret string) *v1AuthHandlers {
	return &v1AuthHandlers{
		userRepo:           userRepo,
		jwtManager:         usecases.NewJWTUseCase(sessionSecret),
		hashPassword:       usecases.NewPassowrdHashUseCase(),
		fastErrorResponses: presentation.NewFastErrorResponses(),
	}
}
