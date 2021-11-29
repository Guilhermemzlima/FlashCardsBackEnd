all: build run

setup:
	go mod vendor
	wire

run:
	go run .

build:
	go mod download
	go mod vendor
	go get -u github.com/google/wire/cmd/wire@v0.5.0
	wire
	CGO_ENABLED=0 GOOS=linux go build -o bin/application


image:
	docker build -t guilhermemzlima/flashcardsbackend:latest .
