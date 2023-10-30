package server

import (
	"github.com/korvised/ilog-shop/modules/item/itemHandlers"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/item/itemRepositories"
	"github.com/korvised/ilog-shop/modules/item/itemUsecases"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"log"
)

func (s *server) itemService() {
	repo := itemRepositories.NewItemRepository(s.db)
	usecase := itemUsecases.NewItemUsecase(repo)
	httpHandler := itemHandlers.NewItemHttpHandler(s.cfg, usecase)
	grpcHandler := itemHandlers.NewItemGrpcHandler(usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)

		itemPb.RegisterItemGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start item gRPC server: %s \n", s.cfg.Grpc.ItemUrl)
		grpcServer.Serve(lis)
	}()

	_ = httpHandler

	router := s.app.Group("/api/v1/item")

	// Health check
	router.GET("", s.healthCheckService)
}
