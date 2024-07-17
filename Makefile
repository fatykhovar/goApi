build:
	@go build -o bin/goApi

run: build
	@./bin/goApi

test:
	@go test -v ./...