
setup_db:
	docker-compose -f docker-compose.db.yml up -d
dev:
	go run main.go ./env/dev/.env.auth

.PHONY: setup_db dev
