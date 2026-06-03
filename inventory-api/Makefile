VERSION ?= v1
CLUSTER  = k8s-lab

build:
	docker build -t inventory-api:$(VERSION) ./backend

load:
	kind load docker-image inventory-api:$(VERSION) --name $(CLUSTER)

deploy:
	kubectl apply -f k8s/

build-load: build load

.PHONY: build load deploy build-load
