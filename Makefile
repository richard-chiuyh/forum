.PHONY: proto build run clean

# 生成 protobuf 代码
proto:
	goctl rpc protoc forum.proto --go_out=./types --go-grpc_out=./types --zrpc_out=. --style goZero -m

# 构建
build:
	go build -o forum forum.go

# 运行
run:
	go run forum.go -f etc/forum.yaml

# 清理
clean:
	rm -f forum
	rm -rf logs/

# 安装依赖
deps:
	go mod download
	go mod tidy

# Docker 构建
docker-build:
	docker build -t forum:latest .

# Docker 运行
docker-run:
	docker run -d --name forum -p 8895:8895 forum:latest
