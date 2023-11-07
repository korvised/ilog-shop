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
	repo := authRepositories.NewAuthRepository(s.db)
	usecase := authUsecases.NewAuthUsecase(s.cfg, repo)
	httpHandler := authHandlers.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := authHandlers.NewAuthGrpcHandler(usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start auth gRPC server: %s \n", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(lis)
	}()

	_ = httpHandler

	router := s.app.Group("/api/v1/auth")

	// Health check
	router.GET("", s.healthCheckService)

	router.POST("/login", httpHandler.Login)
	router.POST("/logout", httpHandler.Logout)
	router.POST("/refresh-token", httpHandler.RefreshToken)
}
