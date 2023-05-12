FROM mysql:8.0

ENV TZ=UTC

COPY ./docker/mysql/initdb.d/init.sql /docker-entrypoint-initdb.d
