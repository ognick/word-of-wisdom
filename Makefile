.PHONY:
.SILENT:
.DEFAULT_GOAL := run

run:
	docker-compose up --build

test:
	go test --count=1  -bench=. -v ./...

lint:
	golangci-lint run