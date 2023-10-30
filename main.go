package main

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/pkg/database"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	// Initialize config
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("error .env path is required")
		}
		return os.Args[1]
	}())

	// Initialize database connection
	db := database.DbConn(ctx, &cfg.Db)
	defer db.Disconnect(ctx)
	log.Println(db)

}
