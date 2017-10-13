#########################
# PROJECT CONFIGURATION #
#########################
SHELL := /bin/bash

SOURCEDIR=.

VERSION=1.0.0

PROJECT_SOURCE=./src/...

#####################
# TASKS DESCRIPTION #
#####################
.PHONY: all lint test $(DOCKER_REDIS_COMMANDS) $(DOCKER_POSTGRES_COMMANDS)

all: lint test

################
# DB MIGRATION #
################

migrate:
	@go run ./src/migrations/migration.go

#####################
# PROTOBUF GENERATE #
#####################

proto-generate:
	@protoc -I$$PWD/src --go_out=plugins=micro:$$PWD/src $$PWD/src/protos/*.proto
	@echo "Success"

###########
# LINTING #
###########

LINTER_CONFIGS=./misc/lint_config.json

lint:
	@gometalinter.v1 --config $(LINTER_CONFIGS) --deadline=60s --vendor $(PROJECT_SOURCE)

#########
# TESTS #
#########

test:
	@go test $(shell go list ./... | grep -v vendor) -cover

#########################
# GENERATE CERTIFICATES #
#########################

generate-certificate:
	@go run $$GOROOT/src/crypto/tls/generate_cert.go --host localhost
	@mv *.pem ./src/

###################
# DOCKER COMMANDS #
###################

docker-init: docker-network-create docker-init-redis docker-init-postgres docker-init-kafka

docker-network-create:
	@docker network create token 2>/dev/null || echo "Network already exists"

###################
# DOCKER STORAGES #
###################

DOCKERFILE=./misc/dockerfiles/Dockerfile

docker-init-%:
	@echo ---------------
	@echo Init $*
	@echo ---------------
	@$(MAKE) docker-build-$*
	@$(MAKE) docker-start-$*

docker-build-%:
	@echo Checking if container exists
	@test -e $(DOCKERFILE).$* 2>/dev/null || (echo Unknown container && exit 1)
	@cat $(DOCKERFILE).$* | docker build -t token-$* -
	@docker rm -f token-$* 2>/dev/null || echo "$* is not running"
#	Retrieve variable with docker additional params from config file
	@. ./misc/docker_config ; \
	DOCKER_PARAMS_VAR=DOCKER_`echo "$*" | tr '[a-z]' '[A-Z]'`; \
	docker run --network token --name token-$* -d $${!DOCKER_PARAMS_VAR} token-$*; \


docker-start-%:
	@docker start token-$*

docker-stop-%:
	@docker stop token-$*
