package server

import (
	"github.com/korvised/ilog-shop/modules/player/playerHandlers"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"github.com/korvised/ilog-shop/modules/player/playerRepositories"
	"github.com/korvised/ilog-shop/modules/player/playerUsecases"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"log"
)

func (s *server) playerService() {
	repo := playerRepositories.NewPlayerRepository(s.db)
	usecase := playerUsecases.NewPlayerUsecase(repo)
	httpHandler := playerHandlers.NewPlayerHttpHandler(s.cfg, usecase)
	grpcHandler := playerHandlers.NewPlayerGrpcHandler(usecase)
	queueHandler := playerHandlers.NewPlayerQueueHandler(s.cfg, usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.PlayerUrl)

		playerPb.RegisterPlayerGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start player gRPC server: %s \n", s.cfg.Grpc.PlayerUrl)
		grpcServer.Serve(lis)
	}()

	_ = httpHandler
	_ = grpcHandler
	_ = queueHandler

	router := s.app.Group("/api/v1/player")

	// Health check
	router.GET("", s.healthCheckService)
}
