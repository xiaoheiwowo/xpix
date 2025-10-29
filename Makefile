.PHONY: build run clean test install

# 项目名称
BINARY_NAME=xpix

# 编译
build:
	go build -o $(BINARY_NAME) main.go

# 运行
run:
	go run main.go

# 清理
clean:
	go clean
	rm -f $(BINARY_NAME)

# 测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod download
	go mod tidy

# 安装到系统
install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/

# 交叉编译
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe main.go

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

