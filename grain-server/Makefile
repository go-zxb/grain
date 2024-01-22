#设置window环境
winEnv:
	go env -w GOOS=windows
	go env -w GO111MODULE=auto
	go env -w GOARCH=amd64
	go env -w CC=gcc
	go env -w  CGO_ENABLED=0

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

mod:
	go mod tidy

# 生成swag文档
swag:
	swag init --generalInfo=./cmd/main.go --output=./docs --parseDependency

gen:
	go run cmd/gen/build.go


android:
	go env -w GOOS=android
	go env -w GO111MODULE=auto
	go env -w GOARCH=arm64
	go env -w  CGO_ENABLED=1

