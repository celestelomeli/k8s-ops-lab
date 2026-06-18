VERSION ?= v1
CLUSTER  = k8s-lab

build:
	docker build -t inventory-api:$(VERSION) ./backend

load:
	kind load docker-image inventory-api:$(VERSION) --name $(CLUSTER)

deploy:
	kubectl apply -f k8s/

build-load: build load

# Unit tests 
test:
	cd backend && go test -v ./...

# Integration tests
test-integration:
	docker compose up -d postgres
	@set -a; . ./.env; set +a; \
		cd backend && DB_HOST=localhost DB_PORT=5432 go test -tags=integration -v ./...; \
		status=$$?; \
		cd .. && docker compose down; \
		exit $$status

.PHONY: build load deploy build-load test test-integration
