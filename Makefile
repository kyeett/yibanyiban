BINARY_NAME=yibanserver
SERVER_PORT?=8080

all: test build
build:
	@go build -o $(BINARY_NAME) -v
test:
	@go test
load-test:
	./load_test.sh
cover:
	@go test -cover
clean:
	@go clean
	rm -f $(BINARY_NAME)
run:
	@go run server.go -port $(SERVER_PORT)