SHELL := /bin/bash
.PHONY: all test-int test-e2e resilence
app:
	@set -a; [ -f .env ] && . ./.env; set +a; \
	docker compose -f docker-compose.base.yaml -f docker-compose.yaml down; \
	docker compose -f docker-compose.base.yaml -f docker-compose.yaml build --build-arg SKIP_LINT=true; \
	docker compose -f docker-compose.base.yaml -f docker-compose.yaml up --scale app=$${APP_REPLICAS} proxy
#need scala project in neighbour directory
all:
	@set -a; [ -f .env ] && . ./.env; set +a; \
	docker compose -f docker-compose.base.yaml -f docker-compose.yaml up -d --scale app=$${APP_REPLICAS} proxy; \
	: > ../go.log; \
	docker compose -f docker-compose.base.yaml -f docker-compose.yaml logs -f proxy > ../go.log 2>&1 & \
	sleep 15; \
	cd ../AsyncFactorial && $(MAKE) service; \
	: > ../scala.log; \
	cd ../AsyncFactorial && docker compose logs -f consumer | grep -E "ERROR|WARN|INFO important" > ../scala.log 2>&1 &
#grep because a large amount of spam

test-int:
	@{ \
		set -e; \
		docker compose -f docker-compose.base.yaml -f docker-compose-test.yaml up -d --build db_test migrate_test app_test kafka1 kafka2 kafka3 kafka-init; \
		echo "waiting for app_test to become healthy..."; \
		for i in $$(seq 1 60); do \
			if docker compose -f docker-compose.base.yaml -f docker-compose-test.yaml exec -T app_test /healthcheck >/dev/null 2>&1; then \
				break; \
			fi; \
			sleep 1; \
		done; \
		docker compose -f docker-compose.base.yaml -f docker-compose-test.yaml exec -T app_test /healthcheck >/dev/null 2>&1; \
		docker compose -f docker-compose.base.yaml -f docker-compose-test.yaml run --rm --no-deps tests; \
		status=$$?; \
		docker compose -f docker-compose.base.yaml -f docker-compose-test.yaml down -v; \
		exit $$status; \
	} > test.log 2>&1

# без кэширования
test-e2e:
	@set -a; [ -f .env ] && . ./.env; set +a; \
	API_BASE_URL=$${API_BASE_URL:-http://localhost:$${SERVER_PORT}} \
		go test ./tests/end-to-end/... -count=1  


#samples: make resilence TEST=app_instance_down_api
#		 make resilence TEST=app_instance_down_factorial
#		make resilence TEST=broker_down
resilence:
	@set -euo pipefail; \
	if [ -z "$${TEST:-}" ]; then echo "Usage: make resilence TEST=<name>"; exit 1; fi; \
	script="./tests/resilence/$${TEST}.sh"; \
	if [ ! -x "$${script}" ]; then echo "Script $${script} not found or not executable"; exit 1; fi; \
	"$${script}"

# НЕ ИСПОЛЬЗОВАТЬ: старые/нестабильные попытки поднять kafka+тесты из make.
# Оставлены как напоминание, чтобы не писать нерабочие конфиги.
test-kafka:
	@echo "test-kafka disabled (deprecated/unstable). Use make test-int or docker compose directly." && exit 1

lint:
	golangci-lint run --config .golangci.yml ./cmd/... ./internal/... ./tests/...

test-load:
	docker rm -f prassign-k6-1 2>/dev/null || true
	docker network prune -f
	docker compose up -d --no-deps k6

fmt:
	gofmt -s -w .

status:
	 docker compose -f docker-compose.base.yaml -f docker-compose.yaml ps

# rebuild: docker compose -f docker-compose.base.yaml -f docker-compose-test.yaml up --build -d app_test

# docker compose down -v --remove-orphans
# 	docker compose build --build-arg SKIP_LINT=true k6
# 	docker compose up -d db app
# 	docker compose up -d --no-deps k6

#docker compose --profile tests up --build k6 - with logs
logs-test-load:
	docker logs -f prassign-k6-1 > out.txt
#docker compose run --rm k6 run /tests/test-create-pr.js


PORTS := 5432 8080 9092 18080 29092 29093 29094

check-ports:
	@echo "Проверка портов: $(PORTS)"
	@blocked=0; \
	for port in $(PORTS); do \
		if lsof -iTCP:$$port -sTCP:LISTEN -t >/dev/null ; then \
			echo "ALERT Порт $$port занят"; \
			blocked=1; \
		else \
			echo " Порт $$port свободен"; \
		fi; \
	done; \
	if [ $$blocked -eq 1 ]; then \
		echo "Некоторые порты заняты, запуск может упасть"; \
		exit 1; \
	else \
		echo "Все порты свободны, можно запускать проект"; \
	fi