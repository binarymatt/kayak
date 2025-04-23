.PHONY: proto
proto:
	buf generate

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: mocks
mocks:
	@mockery 

