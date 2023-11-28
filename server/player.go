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
	repo := playerRepositories.NewPlayerRepository(s.db, s.cfg)
	usecase := playerUsecases.NewPlayerUsecase(repo)
	httpHandler := playerHandlers.NewPlayerHttpHandler(s.cfg, usecase)
	grpcHandler := playerHandlers.NewPlayerGrpcHandler(usecase)
	queueHandler := playerHandlers.NewPlayerQueueHandler(s.cfg, usecase)

	// Kafka
	go queueHandler.DockedPlayerMoney()
	go queueHandler.RollbackPlayerTransaction()

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.PlayerUrl)

		playerPb.RegisterPlayerGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start player gRPC server: %s \n", s.cfg.Grpc.PlayerUrl)
		_ = grpcServer.Serve(lis)
	}()

	router := s.app.Group("/api/v1")

	// Health check
	router.GET("", s.healthCheckService)

	router.POST("/player/register", httpHandler.Register)
	router.POST("/player/add-money", httpHandler.AddPlayerMoney, s.m.Authorization)
	router.GET("/player/profile/:player_id", httpHandler.GetPlayerProfile, s.m.Authorization)
	router.GET("/player/saving-account", httpHandler.GetPlayerSavingAccount, s.m.Authorization)
}
