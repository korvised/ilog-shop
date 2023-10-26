package main

import (
	"fmt"
	"github.com/korvised/ilog-shop/config"
	"log"
	"os"
)

func main() {
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("error .env path is required")
		}
		return os.Args[1]
	}())

	fmt.Println(cfg)
}
