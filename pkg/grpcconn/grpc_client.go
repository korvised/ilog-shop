package grpcconn

import (
	"errors"
	authPb "github.com/korvised/ilog-shop/modules/auth/authPb"
	inventoryPb "github.com/korvised/ilog-shop/modules/inventory/inventoryPb"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	playerPb "github.com/korvised/ilog-shop/modules/player/playerPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func (g *grpcClientFactory) Auth() authPb.AuthGrpcServiceClient {
	return authPb.NewAuthGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Player() playerPb.PlayerGrpcServiceClient {
	return playerPb.NewPlayerGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Item() itemPb.ItemGrpcServiceClient {
	return itemPb.NewItemGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Inventory() inventoryPb.InventoryGrpcServiceClient {
	return inventoryPb.NewInventoryGrpcServiceClient(g.client)
}

func NewGrpcClient(host string) (GrpcClientFactoryHandler, error) {
	opts := make([]grpc.DialOption, 0)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Printf("error: grpc client connection failed: %v", err)
		return nil, errors.New("error: grpc client connection failed")
	}

	return &grpcClientFactory{client: clientConn}, nil
}
