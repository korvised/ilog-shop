package server

import (
	"github.com/korvised/ilog-shop/modules/inventory/inventoryHandlers"
	inventoryPb "github.com/korvised/ilog-shop/modules/inventory/inventoryPb"
	"github.com/korvised/ilog-shop/modules/inventory/inventoryRepositories"
	"github.com/korvised/ilog-shop/modules/inventory/inventoryUsecases"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"log"
)

func (s *server) inventoryService() {
	repo := inventoryRepositories.NewInventoryRepository(s.db, s.cfg)
	usecase := inventoryUsecases.NewInventoryUsecase(s.cfg, repo)
	httpHandler := inventoryHandlers.NewInventoryHttpHandler(s.cfg, usecase)
	grpcHandler := inventoryHandlers.NewInventoryGrpcHandler(usecase)
	queueHandler := inventoryHandlers.NewInventoryQueueHandler(s.cfg, usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.InventoryUrl)

		inventoryPb.RegisterInventoryGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start inventory gRPC server: %s \n", s.cfg.Grpc.InventoryUrl)
		_ = grpcServer.Serve(lis)
	}()

	// Queue handler
	go queueHandler.AddPlayerItem()
	go queueHandler.RollbackAddPlayerItem()

	router := s.app.Group("/api/v1")

	// Health check
	router.GET("", s.healthCheckService)

	router.GET("/inventory/:player_id", httpHandler.GetPlayerItems, s.m.Authorization)
}
