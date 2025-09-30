test.unit:
	cd chat_service && go test -short -v ./...
	cd user_service && go test -short -v ./...

test.integration:
	docker compose -f docker-compose.test.yml up -d --build
	sleep 5
	cd user_service && DB_HOST=localhost DB_PORT=5433 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=dbName \
		go test -tags=integration -v ./tests/integration
	docker compose -f docker-compose.test.yml down -v

linter.notification:
	cd notification-service && golangci-lint run

linter.chat:
	cd chat_service && golangci-lint run

linter.user:
	cd user_service && golangci-lint run

build:
	docker compose up --build

down:
	docker compose down

start:
	docker compose start

stop:
	docker compose stop