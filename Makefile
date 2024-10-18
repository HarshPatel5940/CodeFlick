build:
	@echo "Building application..."
	@go build -o main cmd/main.go

run:
	@go run cmd/main.go

clean:
	@echo "Cleaning..."
	@rm -f main

migrate-up:
	@echo "Migrating up..."
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Migrating down..."
	@go run cmd/migrate/main.go down

lint:
	@golangci-lint run

lint-fix:
	@golangci-lint run --fix

docker-start:
	@echo "Building Docker image..."
	@docker compose up --build -d

docker-restart:
	@echo "Restarting Docker container..."
	@docker compose restart

docker-stop:
	@echo "Stopping Docker container..."
	@docker compose down

docker-logs:
	@echo "Showing Docker logs..."
	@docker compose logs -f

docker-db-up:
	@echo "Starting Mino and Postgres containers..."
	@docker start codeflick-minio codeflick-postgres

docker-db-down:
	@echo "Stopping Mino and Postgres containers..."
	@docker stop codeflick-minio codeflick-postgres

docker-db-restart:
	@echo "Restarting Mino and Postgres containers..."
	@docker restart codeflick-minio codeflick-postgres

docker-delete:
	@echo "Deleting Docker containers..."
	@docker compose down --volumes --networks --remove-orphans

build-templ:
	@templ generate

build-tailwind:
	./tailwindcss -m -i ./public/tailwind.css -o ./public/styles.css $(ARGS)

dev-server:
	@air

dev-templ:
	@templ generate --watch

dev-tailwind:
	@make ARGS="--watch" build-tailwind
