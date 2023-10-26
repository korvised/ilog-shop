package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type (
	Config struct {
		App      App
		Db       Db
		Jwt      Jwt
		Kafka    Kafka
		Grpc     Grpc
		Paginate Paginate
	}

	App struct {
		Name  string
		Url   string
		Stage string
	}

	Db struct{ Url string }

	Jwt struct {
		AccessSecretKey  string
		RefreshSecretKey string
		ApiSecretKey     string
		AccessDuration   int64
		RefreshDuration  int64
	}

	Kafka struct {
		Url    string
		ApiKey string
		Secret string
	}

	Grpc struct {
		AuthUrl      string
		PlayerUrl    string
		ItemUrl      string
		InventoryUrl string
		PaymentUrl   string
	}

	Paginate struct {
		ItemNextPageBaseURL      string
		InventoryNextPageBaseURL string
	}
)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("error loading .env file")
	}

	return Config{
		App: App{
			Name:  os.Getenv("APP_NAME"),
			Url:   os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},
		Db: Db{
			Url: os.Getenv("DB_URL"),
		},
		Jwt: Jwt{
			AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
			RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			ApiSecretKey:     os.Getenv("JWT_API_SECRET_KEY"),
			AccessDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("error loading access duration failed")
				}
				return result
			}(),
			RefreshDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("error loading refresh duration failed")
				}
				return result
			}(),
		},
		Kafka: Kafka{
			Url:    os.Getenv("KAFKA_URL"),
			ApiKey: os.Getenv("KAFKA_API_KEY"),
			Secret: os.Getenv("KAFKA_SECRET"),
		},
		Grpc: Grpc{
			AuthUrl:      os.Getenv("GRPC_AUTH_URL"),
			PlayerUrl:    os.Getenv("GRPC_PLAYER_URL"),
			ItemUrl:      os.Getenv("GRPC_ITEM_URL"),
			InventoryUrl: os.Getenv("GRPC_INVENTORY_URL"),
			PaymentUrl:   os.Getenv("GRPC_PAYMENT_URL"),
		},
		Paginate: Paginate{
			ItemNextPageBaseURL:      os.Getenv("PAGINATE_ITEM_NEXT_PAGE_BASED_URL"),
			InventoryNextPageBaseURL: os.Getenv("PAGINATE_INVENTORY_NEXT_PAGE_BASED_URL"),
		},
	}
}
