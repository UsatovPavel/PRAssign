.PHONY: test-int
app:
	docker compose down
	docker compose  build  --build-arg SKIP_LINT=true
	docker compose up

test-int:
	docker compose -f docker-compose-test.yaml up --build 
	docker compose -f docker-compose-test.yaml down -v

lint:
	golangci-lint run --config .golangci.yml ./... 
test-load:
	docker rm -f prassign-k6-1 2>/dev/null || true
	docker network prune -f
	docker compose up -d --no-deps k6

# docker compose down -v --remove-orphans
# 	docker compose build --build-arg SKIP_LINT=true k6
# 	docker compose up -d db app
# 	docker compose up -d --no-deps k6

#docker compose --profile tests up --build k6 - with logs
logs-test-load:
	docker logs -f prassign-k6-1 > out.txt
#docker compose run --rm k6 run /tests/test-create-pr.js