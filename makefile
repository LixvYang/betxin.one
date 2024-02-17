.PHONY: all build run gotool clean help



MAIN=./cmd/main/main.go

TAG=`date +"%Y%m%d%H%M%S"`
BINARY="betxin_"${TAG}

all: gotool build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY} ${MAIN} 

run:
	CGO_ENABLED=0 go run ./cmd/main/main.go

genup:
	CGO_ENABLED=0 go run main.go gen up -f ./config/config.yaml
	
gendown:
	CGO_ENABLED=0 go run main.go gen down -f ./config/config.yaml

gendata:
	CGO_ENABLED=0 go run main.go gen gendata -f ./config/config.yaml

gen:
	CGO_ENABLED=0 go run main.go gen -f ./config/config.yaml

httpd:
	CGO_ENABLED=0 go run main.go httpd -f ./config/config.yaml

gotool:
	CGO_ENABLED=0 go fmt ./...
	CGO_ENABLED=0 go vet ./...

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm log/*

make push:
	rm log/*
	git add .
	git commit -m "update"
	git push

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
