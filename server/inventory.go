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
	repo := inventoryRepositories.NewInventoryRepository(s.db)
	usecase := inventoryUsecases.NewInventoryUsecase(repo)
	httpHandler := inventoryHandlers.NewInventoryHttpHandler(s.cfg, usecase)
	grpcHandler := inventoryHandlers.NewInventoryGrpcHandler(usecase)
	queueHandler := inventoryHandlers.NewInventoryQueueHandler(s.cfg, usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.InventoryUrl)

		inventoryPb.RegisterInventoryGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start inventory gRPC server: %s \n", s.cfg.Grpc.InventoryUrl)
		grpcServer.Serve(lis)
	}()

	_ = httpHandler
	_ = queueHandler

	router := s.app.Group("/api/v1/inventory")

	// Health check
	router.GET("", s.healthCheckService)
}
