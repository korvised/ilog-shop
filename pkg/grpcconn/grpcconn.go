package grpcconn

import (
	authPb "github.com/korvised/ilog-shop/modules/auth/authPb"
	inventoryPb "github.com/korvised/ilog-shop/modules/inventory/inventoryPb"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"google.golang.org/grpc"
)

type (
	GrpcClientFactoryHandler interface {
		Auth() authPb.AuthGrpcServiceClient
		Player() playerPb.PlayerGrpcServiceClient
		Item() itemPb.ItemGrpcServiceClient
		Inventory() inventoryPb.InventoryGrpcServiceClient
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
	}
)
