.PHONY: proto
proto:
	buf generate

.PHONY: mocks
mocks:
	mockery --all --with-expecter=true --dir=./internal/store

.PHONY: test
test:
	go test -v -cover -coverprofile=coverage.out ./...

.PHONY: local_cluster
local_cluster:
	docker compose up

.PHONY: local_server
local_server:
	go run cmd/kayak/main.go --id 1 --dir data/ --config server1.yaml --host localhost --console

.PHONY: integration_test
integration_test:
	docker compose down
	earthly +docker
	docker compose up -d
	env KAYAK_INTEGRATION_TESTS=true go test -v ./kayak_test.go
	docker compose down

.PHONY: docker
docker:
	docker build . -t kayak:latest

.PHONY: prometheus
prometheus:
	prometheus --config.file=prometheus.yml
