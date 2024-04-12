# ===== RUN =====
.PHONY: build
build:
	docker compose -f docker-compose.yml build api

.PHONY: up
up:
	docker compose -f docker-compose.yml up -d --build db cache api
	docker ps
	make logs

.PHONY: stop
stop:
	docker compose -f docker-compose.yml stop

.PHONY: down
down:
	docker compose -f docker-compose.yml down -v

# ===== LOGS =====
service = api
.PHONY: logs
logs:
	docker compose logs -f $(service)

.PHONY: test-integration
test-integration:
	docker compose -f docker-compose.yml up -d --build test-db
	sleep 2
	go test ./test/integration/...
	#go test -count=50 -bench ./tests/integration/...
	docker compose -f docker-compose.yml down -v test-db
