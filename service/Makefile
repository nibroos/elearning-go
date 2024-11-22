# include .env

#   create_migration:
#     migrate create -ext=sql -dir=internal/database/migrations -seq init

#   migrate_up:
#     migrate -path=internal/database/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

#   migrate_down:
#     migrate -path=internal/database/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down

#   .PHONY: create_migration migrate_up migrate_down

#   build:
#     go build -o bin/main main.go

#   run:
#     go run main.go

# Makefile

# Load environment variables from .env
ifneq (,$(wildcard ./.env))
  include .env
  export
endif

# Variables
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
MIGRATE=migrate -path=internal/database/migrations -database $(DB_URL)
PROTO_DIR=internal/proto
GENERATED_PROTO_DIR=internal/proto/generated
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)
SEEDER_DIR=internal/database/seeders

# Commands

.PHONY: migrate-create
migrate-create:
	$(MIGRATE) create -ext sql -dir internal/database/migrations $(name)

.PHONY: migrate-up
migrate-up:
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE) down 1

.PHONY: migrate-force
migrate-force:
	$(MIGRATE) force $(v-force)
	
.PHONY: seed-create
seed-create:
ifndef name
	$(error name is undefined. Usage: make seed-create name=<seeder_name>)
endif
	TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	echo "Creating seeder file: $(SEEDER_DIR)/$${TIMESTAMP}_$(name).sql"; \
	echo "BEGIN;" > $(SEEDER_DIR)/$${TIMESTAMP}_$(name).sql; \
	echo "INSERT INTO ..." >> $(SEEDER_DIR)/$${TIMESTAMP}_$(name).sql; \
	echo "COMMIT;" >> $(SEEDER_DIR)/$${TIMESTAMP}_$(name).sql

.PHONY: seed-run
seed-run:
ifndef file
	$(error file is undefined. Usage: make seed-run file=<seeder_file>)
endif
	@echo "Running seeder file: $(SEEDER_DIR)/$(file)"
	docker run --rm -v $(PWD)/$(SEEDER_DIR):/seeders --network host postgres:latest \
	psql $(DB_URL) -f /seeders/$(file)

.PHONY: seed-run-prod
seed-run:
ifndef file
	$(error file is undefined. Usage: make seed-run file=<seeder_file>)
endif
	@echo "Running seeder file: $(SEEDER_DIR)/$(file)"
	psql $(DB_URL) -f /seeders/$(file)

.PHONY: build
build:
	go build -o bin/app cmd/main.go

# .PHONY: test
# test:
# 	go test ./...

.PHONY: run
run: build
	./bin/app

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILES)

# Command to initialize the environment
.PHONY: init-env
init-env:
	@echo "Setting up environment variables..."
	@source .env && $(MAKE) migrate-up