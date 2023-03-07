.PHONY:

run_es:
	go run main.go -config=./config/config.yaml

# Set up database.
setup_db:
	docker compose up -d

dev:
	@echo Starting dev docker compose
	docker compose up -d --build
