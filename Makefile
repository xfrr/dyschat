
SERVICES=auth ws-agent rooms
PROTOS=auth rooms
DOCKER_RUN_MODES=dev prod

.PHONY: all
all: build push

.PHONY: $(SERVICES)
build: $(SERVICES)

.PHONY: run
run_$(DOCKER_RUN_MODES):
	@echo "Running $(@:run_%=%) mode"
	docker-compose -f docker/docker-compose.$(@:run_%=%).yml --env-file docker/.$(@:run_%=%).env up

down_$(DOCKER_RUN_MODES):
	@echo "Stopping $(@:down_%=%) mode"
	docker-compose -f docker/docker-compose.$(@:down_%=%).yml --env-file docker/.$(@:down_%=%).env down

.PHONY: proto
proto:
	@echo "Generating gRPC code"
	cd proto && buf generate
	
bench:
	@echo "Running benchmarks"
	go test -bench=. -benchmem ./...