package inventoryRepositories

import (
	"context"
	"errors"
	itemPb "github.com/korvised/ilog-shop/modules/item/itemPb"
	"github.com/korvised/ilog-shop/pkg/grpcconn"
	"github.com/korvised/ilog-shop/pkg/jwtauth"
	"log"
	"time"
)

func (r *inventoryRepository) FindItemInIds(c context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	conn, err := grpcconn.NewGrpcClient(r.cfg.Grpc.ItemUrl)
	if err != nil {
		log.Printf("Error: gRPC client connection failed: %v \n", err)
		return nil, errors.New("error: gRPC client connection failed")
	}

	jwtauth.SetApiKeyInContext(&ctx)
	result, err := conn.Item().FindItemInIds(ctx, req)
	if err != nil {
		log.Printf("Error: FindItemInIds: %v \n", err)
		return nil, errors.New("error: item not found")
	}

	if result == nil || result.Items == nil || len(result.Items) == 0 {
		log.Println("Error FindItemInIds: result is nil or empty")
		return nil, errors.New("error: item not found")
	}

	return result, nil
}
