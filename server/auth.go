package server

import (
	"github.com/korvised/ilog-shop/modules/auth/authHandlers"
	"github.com/korvised/ilog-shop/modules/auth/authRepositories"
	"github.com/korvised/ilog-shop/modules/auth/authUsecases"
)

func (s *server) authService() {
	repo := authRepositories.NewAuthRepository(s.db)
	usecase := authUsecases.NewAuthUsecase(repo)
	httpHandler := authHandlers.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := authHandlers.NewAuthGrpcHandler(usecase)

	_ = httpHandler
	_ = grpcHandler

	router := s.app.Group("/api/v1/auth")

	// Health check
	router.GET("", s.healthCheckService)
}
