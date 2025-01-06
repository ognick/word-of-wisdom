.PHONY:
.SILENT:
.DEFAULT_GOAL := run

wire:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

run:
	go run -mod=mod github.com/google/wire/cmd/wire ./...
	docker compose up --build

test:
	go test --count=1  -bench=. -v ./...

lint:
	golangci-lint run --fix --verbose