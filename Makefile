.PHONY: proto
proto:
	buf generate

.PHONY: test
test:
	go test -cover ./...

.PHONY: coverage
coverage:
	go test -v ./... -coverprofile=coverage.out
	@cat coverage.out | grep -v "mocks.go" | grep -v "gen" > cover.out
	go tool cover -html=cover.out
	@rm cover.out 
	@rm coverage.out
.PHONY: all_tests
all_tests:
	env RUN_LONG_TESTS="" go test -v ./...

.PHONY: mocks
mocks:
	@mockery 

.PHONY: run
run:
	go run cmd/kayak/main.go

.PHONY: grpc/protoset
grpc/protoset:
	@buf build -o k6/kayak.protoset --path proto/kayak

grpc/ui:
	grpcui -plaintext -protoset ./k6/kayak.protoset localhost:8080

.PHONY: load
load:
	k6 run k6/loadtest.js

server1:
	go run cmd/kayak/main.go --node_id=server1 --listen_address=localhost:8080 --raft_address=localhost:1200 --data_dir=./data/server1 --raft_data_dir=./raft_data/raft_data1
server2:
	go run cmd/kayak/main.go --node_id=server2 --listen_address=localhost:8081 --raft_address=localhost:1201 --data_dir=./data/server2 --raft_data_dir=./raft_data/raft_data2 --join_addr=localhost:8080
server3:
	go run cmd/kayak/main.go --node_id=server3 --listen_address=localhost:8082 --raft_address=localhost:1202 --data_dir=./data/server3 --raft_data_dir=./raft_data/raft_data3 --join_addr=localhost:8080
