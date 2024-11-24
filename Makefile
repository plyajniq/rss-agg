.PHONY: docs build start

docs:
	swag init -g cmd/main.go

build: docs
	go build -o app cmd/main.go

start: build
	./app
