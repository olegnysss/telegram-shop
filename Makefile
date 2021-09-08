.PHONY:

build:
	go build -o ./.bin/bot src/bot/main.go

run: build
	./.bin/bot