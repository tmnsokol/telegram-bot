.PHONY:

all: build build-image start-container

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t telegram-bot:0.2 .

start-container:
	docker run --env-file .env -p 80:80 telegram-bot:0.2