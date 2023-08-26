
#设置window环境
.PHONY: winEnv
winEnv:
	go env -w GOOS=windows
	go env -w GO111MODULE=auto
	go env -w GOARCH=amd64
	go env -w CC=gcc
	go env -w  CGO_ENABLED=0

.PHONY: build
build:
# Linux arm64版本
	go env -w  GOOS=linux
	go env -w  GOARCH=arm64
	go env -w  CGO_ENABLED=0
	go env -w GO111MODULE=auto
	go build  -ldflags="-s -w"  -o ./build/linux/grain_linux_arm64 ./cmd/main.go
# Linux amd64版本
	go env -w  GOARCH=amd64
	go build   -ldflags="-s -w" -o ./build/linux/grain_linux_amd64 ./cmd/main.go
# Windows 版本
	go env -w  GOARCH=amd64
	go env -w  GOOS=windows
	go build   -ldflags="-s -w" -o ./build/windows/grain_win_amd64.exe ./cmd/main.go

.PHONY: mod
mod:
	go mod tidy

# 生成swag文档
.PHONY: swag
swag:
	swag init --generalInfo=./cmd/main.go --output=./docs --parseDependency

.PHONY: gen
gen:
	go run cmd/gen/build.go