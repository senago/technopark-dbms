# path to actual config - the one that is copied to the docker container
CONFIG:=resources/config/config.yaml

# path to docker compose file
DCOMPOSE:=docker-compose.yaml

DCOMPOSE_BUILD_ARGS=--build-arg CONFIG=${CONFIG} --parallel

# path to external config which will copied to CONFIG
CONFIG_PATH=resources/config/config_default.yaml

# improve build time
DOCKER_BUILD_KIT:=COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1

all: down build up

debug: down build up-debug

down:
	docker-compose -f ${DCOMPOSE} down --remove-orphans

build:
	cp ${CONFIG_PATH} ${CONFIG}
	${DOCKER_BUILD_KIT} docker-compose build ${DCOMPOSE_BUILD_ARGS}

up:
	docker-compose --compatibility -f ${DCOMPOSE} up -d --remove-orphans

up-debug:
	docker-compose --compatibility -f ${DCOMPOSE} up --remove-orphans

# Vendoring is useful for local debugging since you don't have to
# reinstall all packages again and again in docker
mod:
	go mod tidy -compat=1.18 && go mod vendor && go install ./...
