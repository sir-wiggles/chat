# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

DOCKER_NETWOR=chat_default
MIGRATION_DIR=/migrations
DATABASE=postgres://admin:admin@postgres:5432/chat?sslmode=disable
POSTGRES_SERVICE_NAME=postgres
CASSANDRA_SERVICE_NAME=cassandra

SERVICES = api web
DOCKER_BUILD_TARGETS = $(foreach service, $(SERVICES), docker-build-$(service))

PROJECT_NAME=chat

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) .

.PHONY: docker-build $(DOCKER_BUILD_TARGETS)
docker-build: $(DOCKER_BUILD_TARGETS)

$(DOCKER_BUILD_TARGETS):
	docker-compose build $(subst docker-build-,,$@)


up:
	docker-compose up -d postgres
	docker-compose up api web

api-up:
	docker-compose up api

web-up:
	docker-compose up web

psql:
ifeq ($(shell docker-compose ps | grep $(PROJECT_NAME)_$(POSTGRES_SERVICE_NAME) | wc -l), 0)
	docker-compose up -d $(POSTGRES_SERVICE_NAME)
endif
	docker-compose exec $(POSTGRES_SERVICE_NAME) psql -U admin -h localhost -d chat

cqlsh:
ifeq ($(shell docker-compose ps | grep $(PROJECT_NAME)_$(CASSANDRA_SERVICE_NAME) | wc -l), 0)
	docker-compose up -d $(CASSANDRA_SERVICE_NAME)
endif
	docker-compose exec $(CASSANDRA_SERVICE_NAME) cqlsh

build-web:
	rm -rf ./api/static
	cd web && yarn build
	mv ./web/dist ./api/static/

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

migrate-up:
	docker run -v $(PWD)$(MIGRATION_DIR):/migrations --network chat_default migrate/migrate -path='$(MIGRATION_DIR)/' -database='$(DATABASE)' up $(count)

migrate-down:
	docker run -v $(PWD)$(MIGRATION_DIR):/migrations --network chat_default migrate/migrate -path='$(MIGRATION_DIR)/' -database='$(DATABASE)' down $(count)

migrate-force:
	docker run -v $(PWD)$(MIGRATION_DIR):/migrations --network chat_default migrate/migrate -path='$(MIGRATION_DIR)/' -database='$(DATABASE)' force $(version)

go-get:
	docker-compose run api go get -v ./...
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
