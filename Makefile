ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ROOT=$(shell pwd)

run:
	go run main.go

build: build-darwin build-windows build-linux

build-darwin:
	source .env && cd build; sh build-darwin.sh

build-windows:
	source .env && cd build; sh build-windows.sh

build-linux: docker-build-linux docker-clean

docker-build-linux:
	docker build -t crypto-auto -f build/linux/Dockerfile .
	docker run -v $(ROOT)/bin:/crypto-auto/bin -t crypto-auto bash -c 'export VERSION=${VERSION} && export NAME=${NAME} && export NAME_LOWER=${NAME_LOWER} && cd build; bash build-linux.sh'

docker-clean:
	docker rm $(shell docker ps --all -q)
	docker rmi $(shell docker images | grep crypto-auto | tr -s ' ' | cut -d ' ' -f 3)
