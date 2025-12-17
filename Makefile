all: build up check

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

check:
	@echo "Checking services..."
	@sleep 5
	@curl -s http://localhost:8080/api/contracts || echo "Server might not be ready"
	@echo "\n"
	@curl -s http://localhost:3000/ask -X POST -H "Content-Type: application/json" -d '{"question":"ping"}' || echo "AI Agent might not be ready"
	@echo "\nDone."
