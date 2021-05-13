.PHONY: build
build:
	go build -v ./cmd/main.go


.PHONY: run
run:
	go run ./cmd/main.go


.DEFAULT_GOAL := run