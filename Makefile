all: build run

setup:
	go mod vendor
	wire

build:
	go build -o bin/FlashCardsBackEnd

run:
	go run .

