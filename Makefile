# Simple Makefile for a Go project
# Build the application
all: build

build:
	@echo "Building..."

	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

migrate-up:
	@echo "Migrating up..."
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Migrating down..."
	@go run cmd/migrate/main.go down

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
