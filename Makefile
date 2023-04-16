
# Dependend on where the Makefile is used, in a docker environment or not,
# the appropriate config file is loaded.
ifneq ($(wildcard /.dockerenv),)
  export QUICK_CONFIGFILE=config-dev-docker.toml
else
  export QUICK_CONFIGFILE=config-dev-local.toml
endif

CURRENT_UID := $(shell id -u)
CURRENT_GID := $(shell id -g)

export CURRENT_UID
export CURRENT_GID

run:
	docker-compose up -d

down:
	docker-compose down

build-server:
	go build -o main cmd/main.go

local-server:
	@echo Using '$(QUICK_CONFIGFILE)' configuration file
	go run cmd/main.go

watch:
	reflex --config=".reflex.conf" --decoration="none"

tidy: 
	go mod tidy

logs:
	docker-compose logs -f

test: export QUICK_CONFIGFILE=config-test.toml
test:
	go test -v -cover ./...

docs:
	swag fmt && swag init  -g ./cmd/main.go


.PHONY: run down build-server tidy local-server watch logs test docs
