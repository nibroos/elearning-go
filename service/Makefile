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
	@echo "Creating migration file: $(MIGRATE_DIR)/$(name)"
	docker compose -f ../docker/docker-compose-dev.yml run --rm \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-e POSTGRES_PORT=$(POSTGRES_PORT) \
	service \
	$(MIGRATE) create -ext sql -dir internal/database/migrations $(name)

.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up inside Docker container"
	docker compose -f ../docker/docker-compose-dev.yml run --rm \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-e POSTGRES_PORT=$(POSTGRES_PORT) \
	service \
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down:
	@echo "Running migrations down inside Docker container"
	docker compose -f ../docker/docker-compose-dev.yml run --rm \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-e POSTGRES_PORT=$(POSTGRES_PORT) \
	service \
	$(MIGRATE) down 1

.PHONY: migrate-force
migrate-force:
	@echo "Running migrations force inside Docker container"
	docker compose -f ../docker/docker-compose-dev.yml run --rm \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-e POSTGRES_PORT=$(POSTGRES_PORT) \
	service \
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
	@echo "Running seeder file: $(SEEDER_DIR)/$(file) inside Docker container"
	docker compose -f ../docker/docker-compose-dev.yml run --rm \
	-v $(PWD)/$(SEEDER_DIR):/seeders \
	service \
	psql postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@user-postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable -f internal/database/seeders/$(file)

.PHONY: build
build:
	go build -o bin/app cmd/main.go

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