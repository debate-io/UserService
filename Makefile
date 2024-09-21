SEP="========================================================"
GQLGEN ?= github.com/99designs/gqlgen
REGISTRY_HOST ?= registry.owao.space
SERVICE_NAME ?= service-auth

################################################################################################################

.PHONY: create-empty-db
create-empty-db:
	$(call _info, $(SEP))
	$(call _info,"Creating empty db")
	$(call _info, $(SEP))
	docker run -d --name auth_db -p 5432:5432 -e POSTGRES_PASSWORD=123123 -e POSTGRES_USER=postgres -e POSTGRES_DB=auth_db postgres:alpine

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
SERVICE_NAME=service-auth
POSTGRES_DSN=postgresql://postgres:123123@localhost:5432/auth_db?sslmode=disable
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