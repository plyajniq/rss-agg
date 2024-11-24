.PHONY: docs build start

docs:
	swag init -g cmd/main.go -o ./docs

build: docs
	go build -o bin/rssagg cmd/main.go

start: build
	./bin/rssagg
