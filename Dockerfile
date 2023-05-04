FROM golang:1.19-alpine

# Set up the working directory (The directory does not exist and will be created automatically)
WORKDIR /var/www/app/go-todolist-aws

RUN go install github.com/cosmtrek/air@v1.40.4; \
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2; \
    go install github.com/swaggo/swag/cmd/swag@v1.6.0

# Copy the local file/directory to the specified location in the image file
COPY ./go.mod ./go.sum ./

RUN go mod download

RUN air init

CMD ["air", "-c", ".air.toml"]

EXPOSE 9753
