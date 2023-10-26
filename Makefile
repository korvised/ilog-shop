
setup_db:
	docker-compose -f docker-compose.db.yml up -d

.PHONY: setup_db
