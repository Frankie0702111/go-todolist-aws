include .env

DOCKER = docker compose exec server
number :=

# General DB structure
migrate-up:
	$(DOCKER) migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}" -path ./migrations up $(number)

migrate-down:
	$(DOCKER) migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}" -path ./migrations down $(number)

# Test DB structure
migrate-test-up:
	$(DOCKER) migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}_test" -path ./migrations up $(number)
migrate-test-down:
	$(DOCKER) migrate -database "${DB}://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}_test" -path ./migrations down $(number)

# Unit test
go-test:
	make migrate-test-up;
	go test -v repository/authRepository/authRepository_test.go;
	make migrate-test-down;

# Compilation project
build:
	GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
	zip -D -j -r bootstrap.zip bootstrap
