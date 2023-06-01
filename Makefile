include .env

DOCKER = docker compose exec server
YMD = _$$(date +'%Y%m%d')
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
	@set -e; \
	make migrate-test-up; \
	$(DOCKER) go test -v repository/authRepository/authRepository_test.go -json > ./log/authRepository_test$(YMD).log; \
	$(DOCKER) go test -v service/authService/authService_test.go -json > ./log/authService_test$(YMD).log; \
	$(DOCKER) go test -v service/jwtService/jwtService_test.go -json > ./log/jwtService_test$(YMD).log; \
	$(DOCKER) go test -v controller/auth/authController_test.go -json > ./log/authController_test$(YMD).log; \
	make migrate-test-down;

# Generate swagger
swagger:
	$(DOCKER) swag init

# Compilation project
build:
	GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
	zip -D -j -r bootstrap.zip bootstrap
