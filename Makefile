# import config.
# You can change the default config with `make cnf="config_special.env" build`
cnf ?= dev.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

# TODO grep the version from version file
VERSION=0.1

# HELP
# This will output the help for each task
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

run: ## run go run for project
	source dev.env && go run main.go

build: ## run go build for project
	go build main.go

version: ## Output the current version
	@echo $(VERSION)

lint: ## Lint the source code
	go vet github.com/tonydonlon/eventservice

watch: ## start service in dev mode; watch for file changes and bounce server
	reflex -r '\.go$' go run main.go

# DOCKER TASKS
# Build the container
docker-build: ## Build the container
	docker build -t $(APP_NAME) .

docker-run: ## Run container configured in `dev.env`
	docker run -i -t --rm --env-file=./dev.env --name="$(APP_NAME)" $(APP_NAME)

docker-stop: ## Stop and remove a running container
	docker stop $(APP_NAME); docker rm $(APP_NAME)

# TODO this does not work on linux host `--add-host=host.docker.internal:host-gateway`
docker-shell: ## Run container and login to bash shell
	docker run -it --rm --env-file=./dev.env --name="$(APP_NAME)" $(APP_NAME) /bin/sh
