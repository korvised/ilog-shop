package server

import (
	"github.com/korvised/ilog-shop/modules/item/itemHandlers"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/modules/item/itemRepositories"
	"github.com/korvised/ilog-shop/modules/item/itemUsecases"
	"github.com/korvised/ilog-shop/modules/middleware"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"log"
)

func (s *server) itemService() {
	repo := itemRepositories.NewItemRepository(s.db)
	usecase := itemUsecases.NewItemUsecase(s.cfg, repo)
	httpHandler := itemHandlers.NewItemHttpHandler(s.cfg, usecase)
	grpcHandler := itemHandlers.NewItemGrpcHandler(usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpcconn.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)

		itemPb.RegisterItemGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Start item gRPC server: %s \n", s.cfg.Grpc.ItemUrl)
		_ = grpcServer.Serve(lis)
	}()

	_ = httpHandler

	router := s.app.Group("/api/v1")

	// Health check
	router.GET("", s.healthCheckService)

	router.POST("/item", httpHandler.CreateItem, s.m.Authorization, s.m.Roles(middleware.RoleAdmin))
	router.GET("/item", httpHandler.GetItems, s.m.Authorization)
	router.GET("/item/:item_id", httpHandler.GetItem, s.m.Authorization)
	router.PUT("/item/:item_id", httpHandler.UpdateItem, s.m.Authorization, s.m.Roles(middleware.RoleAdmin))
	router.PATCH("/item/:item_id/enable", httpHandler.EnableItem, s.m.Authorization, s.m.Roles(middleware.RoleAdmin))
	router.PATCH("/item/:item_id/disable", httpHandler.DisableItem, s.m.Authorization, s.m.Roles(middleware.RoleAdmin))
}
