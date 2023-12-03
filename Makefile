DC = docker-compose
COMPOSE_FILE = docker-compose.yaml

.PHONY: build, run, drop-all, all, logs

build:
	go build -o ./.bin/app cmd/app/main.go
run: build
	./.bin/app
lint:
	golangci-lint run ./... --config=./.golangci.yml
lint-fast:
	golangci-lint run ./... --fast --config=./.golangci.yml
all:
	${DC} -f ${COMPOSE_FILE} up --build -d
drop-all:
	${DC} -f ${COMPOSE_FILE} down
logs:
	${DC} -f ${COMPOSE_FILE} logs -f