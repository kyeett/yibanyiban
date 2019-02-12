
BINARY_NAME=yibanserver

all: test build
build:
	go build -o $(BINARY_NAME) -v
test:
	go test
cover:
	go test -cover
clean:
	go clean
	rm -f $(BINARY_NAME)
run:
	echo "yay"