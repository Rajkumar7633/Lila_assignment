.PHONY: backend frontend dev ingest docker-build

backend:
	cd backend && CGO_ENABLED=1 DATA_PATH=../data/player_data go run ./cmd/server

frontend:
	cd frontend && npm run dev

dev:
	@echo "Run 'make backend' and 'make frontend' in separate terminals"

ingest:
	DATA_PATH=data/player_data python3 scripts/ingest.py

docker-build:
	docker compose build backend

docker-up:
	docker compose up --build
