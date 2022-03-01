SHELL=/bin/zsh

all: init

init:
	git submodule update --init --recursive
	cd external/minio && go mod tidy
	cd external/minio-go && go mod tidy

up:
	docker-compose -f docker-compose/docker-compose.yml up -d

down:
	docker-compose -f docker-compose/docker-compose.yml down

files:
	go test tools/files_test.go