.PHONY: all build run up deps docs dev test test-cov format

GOPATH ?= $(HOME)/go

all: build

build:
	go build -o ms-data-crawler main.go

run:
	go run -race main.go

seed:
	go run main.go -s

seed-docker:
	docker-compose exec app go run main.go -s

up:
	docker-compose up -d --force-recreate

up-mongo:
	docker-compose up -d --force-recreate mongo mongo-admin

mod:
	go mod download

deps: mod
	( cd /tmp; \
		go get \
			github.com/cespare/reflex \
			github.com/swaggo/swag/cmd/swag \
			 github.com/golang/mock/mockgen \
		)

docs:
	$(GOPATH)/bin/swag init

dev:
	DC_APP_ENV=dev $(GOPATH)/bin/reflex -s -r '\.go$$' make format run

dev2:
	DC_APP_ENV=test DC_HTTP_ADDRESS=0.0.0.0:9004 $(GOPATH)/bin/reflex -s -r '\.go$$' make format run

test:
	go test ./core/...

test-cov:
	    go test -coverprofile=cover.txt ./core/... && go tool cover -html=cover.txt -o cover.html

format:
	go fmt ./...

generate-mock:
	go generate -v ./...
