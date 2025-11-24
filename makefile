.PHONY: test-int
app:
	docker compose down
	docker compose build
	docker compose up

test-int:
	docker compose -f docker-compose-test.yaml up --build 
	docker compose -f docker-compose-test.yaml down -v

lint:
	golangci-lint run --config .golangci.yml ./...
