SEP="========================================================"
GQLGEN ?= github.com/99designs/gqlgen
REGISTRY_HOST ?= registry.owao.space
SERVICE_NAME ?= service-auth

################################################################################################################

.PHONY: create-empty-db run

# Команда для создания базы данных, если контейнер не существует
create-empty-db:
	$(call _info, $(SEP))
	$(call _info,"Checking if the database container exists")
	$(call _info, $(SEP))
	@if [ $$(docker ps -a -f name=auth_db --format '{{.Names}}') = "auth_db" ]; then \
		echo "Database container 'auth_db' already exists."; \
		if [ $$(docker inspect -f '{{.State.Running}}' auth_db) = "false" ]; then \
			echo "Starting existing container 'auth_db'..."; \
			docker start auth_db; \
		else \
			echo "Database container 'auth_db' is already running."; \
		fi \
	else \
		echo "Database container 'auth_db' does not exist. Creating a new one..."; \
		docker run -d --name auth_db -p 5434:5432 -e POSTGRES_PASSWORD=123123 -e POSTGRES_USER=user-owner -e POSTGRES_DB=user-db postgres:alpine; \
	fi

# Команда для запуска приложения
run: create-empty-db
	$(call _info, $(SEP))
	$(call _info,"Starting the application")
	$(call _info, $(SEP))
	go run cmd/app/main.go

################################################################################################################

.PHONY: install-generator
install-generator:
	go get $(GQLGEN)

.PHONY: gen
gen:
	go get $(GQLGEN) && go run $(GQLGEN) generate --config ./internal/interface/graphql/gqlgen.yml

################################################################################################################

.PHONY: lint-code
lint-code:
	golangci-lint run
################################################################################################################
.PHONY: env

define ENV_SAMPLE
SERVICE_NAME=user-service
POSTGRES_DSN=postgresql://user-owner:123123@localhost:5434/user-db?sslmode=disable
SERVER_ADDRESS=:9090
IS_DEBUG=true
JWT_SECRET_AUTH=nigganigga
JWT_SECRET_MESSAGES=chungachanga
DAYS_AUTH_EXPIRES=31
DAYS_RECOVERY_EXPIRES=30
endef
export ENV_SAMPLE
env:
	@if [ ! -f ".env" ];\
		then echo "$$ENV_SAMPLE" > .env;\
	 fi

################################################################################################################
.PHONY: deploy
deploy:
	docker build -t $(REGISTRY_HOST)/$(SERVICE_NAME) . && docker push $(REGISTRY_HOST)/$(SERVICE_NAME)