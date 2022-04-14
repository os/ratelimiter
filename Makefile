.PHONY: run run-local test test-local

run:
	docker-compose up --build --abort-on-container-exit

run-local:
	go run main.go

test:
	docker-compose -f testing.docker-compose.yml up --build --abort-on-container-exit

test-local:
	go test -v ./...
