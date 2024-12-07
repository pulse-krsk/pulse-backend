rm: 
	docker compose stop \
	&& docker compose down -v \
	&& docker compose rm \
	&& docker image rm ht-backend:local \
	&& sudo rm -rf pgdata/

up: 
	docker build -t ht-backend:local . \
	&& docker compose up --detach --force-recreate

up-db:
	docker compose up -d postgresql

rm-db:
	docker compose stop \
	&& docker compose rm \
	&& sudo rm -rf pgdata/

# Остановка и удаление всех сервисов, очистка данных и удаление образа
function rm {
    docker compose stop `
    && docker compose down -v `
    && docker compose rm `
    && docker image rm pulse-backend:local `
    && Remove-Item -Recurse -Force pgdata
}

# Сборка образа и запуск всех сервисов
function up {
    docker build -t pulse-backend:local . `
    && docker compose up --detach --force-recreate
}

# Запуск только базы данных
function up-db {
    docker compose up -d postgresql
}

# Остановка и удаление базы данных, очистка данных
function rm-db {
    docker compose stop `
    && docker compose rm `
    && Remove-Item -Recurse -Force pgdata
}
