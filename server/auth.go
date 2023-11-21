package server

import (
	"github.com/korvised/ilog-shop/modules/auth/authHandlers"
	authPb "github.com/korvised/ilog-shop/modules/auth/authPb"
	"github.com/korvised/ilog-shop/modules/auth/authRepositories"
	"github.com/korvised/ilog-shop/modules/auth/authUsecases"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"log"
)

func (s *server) authService() {
	repo := authRepositories.NewAuthRepository(s.db, s.cfg)
	usecase := authUsecases.NewAuthUsecase(s.cfg, repo)
	httpHandler := authHandlers.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := authHandlers.NewAuthGrpcHandler(usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start auth gRPC server: %s \n", s.cfg.Grpc.AuthUrl)
		_ = grpcServer.Serve(lis)
	}()

	router := s.app.Group("/api/v1")

	// Health check
	router.GET("", s.healthCheckService)

	router.POST("/auth/login", httpHandler.Login)
	router.POST("/auth/logout", httpHandler.Logout, s.m.Authorization)
	router.POST("/auth/refresh-token", httpHandler.RefreshToken, s.m.Authorization)
}
