.PHONY: docs build start

docs:
	swag init -g cmd/main.go -o ./docs

build-linux: docs
	GOOS=linux GOARCH=amd64 go build -o ./bin/rssagg cmd/main.go

build-local: docs
	go build -o ./bin/rssagg cmd/main.go

start: build-local
	./bin/rssagg
