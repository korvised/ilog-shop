
setup_db:
	docker-compose -f docker-compose.db.yml up -d
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

.PHONY: setup_db start_auth start_inventory start_item start_payment start_player
