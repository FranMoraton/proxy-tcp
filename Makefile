UID=$(shell id -u)
GID=$(shell id -g)
APP_SERVICE=myapp
TCP_SERVER_MOCK=tcp_server
DUMMY_API=dummy_api

init: erase cache-folders build git-hooks start

erase:
		docker compose down -v

build:
		docker compose build --no-cache && \
		docker compose pull
start:
		docker compose up -d

stop:
		docker compose stop

bash:
		docker compose run --rm -u ${UID}:${GID} ${APP_SERVICE} sh

logs:
		docker compose logs -f ${APP_SERVICE}

git-hooks:
		cp -r ./git_hooks/* .git/hooks

tcp-server-logs:
		docker compose logs -f ${TCP_SERVER_MOCK}

api-dummy-logs:
		docker compose logs -f ${DUMMY_API}

cache-folders:
		mkdir -p ~/.go-cache && chown ${UID}:${GID} ~/.go-cache

curl:
		curl -X POST http://localhost:8080/send
