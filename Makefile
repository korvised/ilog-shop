
setup_db:
	docker-compose -f docker-compose.db.yml up -d
gen_auth_grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./modules/auth/authPb/authPb.proto
start_auth:
	go run main.go ./env/dev/.env.auth
start_inventory:
	go run main.go ./env/dev/.env.inventory
start_item:
	go run main.go ./env/dev/.env.item
start_payment:
	go run main.go ./env/dev/.env.payment
start_player:
	go run main.go ./env/dev/.env.player

.PHONY: setup_db gen_auth_grpc start_auth start_inventory start_item start_payment start_player
