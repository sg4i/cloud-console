 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=cloudconsole
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd
PROTO_PATH=./proto
DOCKERFILE_PATH=deployment/Dockerfile


all: tidy test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run: build
	./$(BINARY_NAME)

deps:
	$(GOGET) -v -d ./...


# 添加新的目标用于生成 protobuf 代码
proto:
	protoc -I. -Ithird_party \
	    --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		$(PROTO_PATH)/*.proto

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

docker-build:
	docker build -t $(BINARY_NAME):latest -f $(DOCKERFILE_PATH) .

tidy:
	$(GOCMD) mod tidy

.PHONY: all build test clean run deps build-linux docker-build tidy proto