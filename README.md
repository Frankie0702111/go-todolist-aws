# Project description
This project will use Gin, GORM, MySQL (RDS), Redis (EC2), api-gateway, Lambda, S3, VPC, IAM to complete a simple to-do list.<br>
**I am sorry that this project only shows the code and cannot teach how to set up AWS**

# Contents
 - [Software requirements](#software-requirements)
 - [Project plugins](#project-plugins)
 - [How to build project](#how-to-build-project)
 - [Folder structure](#folder-structure)
 - [Folder definition](#folder-definition)

# Software requirements
 - **Compilation tools**
    - [Vscode](https://code.visualstudio.com/)
 - **Database**
    - [MySQL](https://aws.amazon.com/tw/rds/): v8
    - [Redis](https://aws.amazon.com/tw/ec2/): v6
 - **Programming language**
    - [Go](https://go.dev/dl/): v1.19
 - **Deveops**
    - [Docker GUI](https://www.docker.com/products/docker-desktop/)
 - **Other**
    - [Postman](https://www.postman.com/downloads/)

# Project plugins
- [AWS sdk for go v2](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
- [AWS api-gateway](https://github.com/apex/gateway)
- [Crypto](https://pkg.go.dev/golang.org/x/crypto)
- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [Go-redis](https://github.com/go-redis/redis)
- [Golang-jwt](https://github.com/golang-jwt/jwt)
- [MySQL](https://github.com/go-gorm/mysql)
- [smapping](https://github.com/mashingan/smapping)
- [swagger](https://github.com/swaggo/swag)
- [testify](https://github.com/stretchr/testify)
- [uuid](https://github.com/gofrs/uuid)

# How to build project
## 1.Clone GitHub project to local
```bash
git clone https://github.com/Frankie0702111/go-todolist-aws.git
```

## 2.Generate config and set up environment
```bash
cd go-todolist-aws/config
cp config.go.example config.go

# Set up basic information, such as database, AWS, JWT
vim config.go
```

## 3.Build docker image and start
```bash
cd ../go-todolist-aws

# Create docker image
docker compose build --no-cache

# Run docker
docker compose up -d

# Stop docker
docker compose stop
```

## 4.Generate db migrations
```bash
# Up all migration
make migrate-up

# Down all migration
make migrate-down

# Specify batch up or down (If you want to go down to a specific file, it is recommended to open a new folder)
make migrate-up number=1
make migrate-down number=1
```

## 5.Build Swagger API documentation
```bash
make swagger
```
- [Swagger API doc](http://localhost:9753/api/swagger/index.html)

## 6.Test authentication example
```bash
make go-test
```

## 7.Build golang project for aws lambda and compress to zip
```bash
make build
```

# Folder structure
```
├── LICENSE
├── Makefile
├── README.md
├── config
│   └── config.go.example
├── controller
│   ├── auth
│   │   ├── authController.go
│   │   └── authController_test.go
│   ├── category
│   │   └── categoryController.go
│   └── task
│       └── taskController.go
├── docker
│   ├── golang
│   │   └── Dockerfile
│   └── mysql
│       ├── Dockerfile
│       └── initdb.d
│           └── init.sql
├── docker-compose.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── log
├── main.go
├── middleware
│   ├── cors.go
│   └── jwt.go
├── migrations
│   ├── 20221129000000_create_users_table.down.sql
│   ├── 20221129000000_create_users_table.up.sql
│   ├── 20221129000001_create_categories_table.down.sql
│   ├── 20221129000001_create_categories_table.up.sql
│   ├── 20221129000002_create_tasks_table.down.sql
│   └── 20221129000002_create_tasks_table.up.sql
├── model
│   ├── category.go
│   ├── task.go
│   └── user.go
├── repository
│   ├── authRepository
│   │   ├── authRepository.go
│   │   ├── authRepository_test.go
│   │   └── repository.go
│   ├── categoryRepository
│   │   ├── categoryRepository.go
│   │   └── repository.go
│   ├── redisRepository
│   │   ├── redisRepository.go
│   │   └── repository.go
│   ├── s3Repository
│   │   ├── repository.go
│   │   └── s3Repository.go
│   └── taskRepository
│       ├── repository.go
│       └── taskRepository.go
├── request
│   ├── authRequest
│   │   └── authRequest.go
│   ├── categoryRequest
│   │   └── categoryRequest.go
│   ├── publicRequest.go
│   └── taskRequest
│       └── taskRequest.go
├── router
│   ├── authRouter
│   │   └── authRouter.go
│   ├── categoryRouter
│   │   └── categoryRouter.go
│   ├── router.go
│   └── taskRouter
│       └── taskRouter.go
├── service
│   ├── authService
│   │   ├── authService.go
│   │   ├── authService_test.go
│   │   └── service.go
│   ├── categoryService
│   │   ├── categoryService.go
│   │   └── service.go
│   ├── jwtService
│   │   ├── jwtService.go
│   │   ├── jwtService_test.go
│   │   └── service.go
│   └── taskService
│       ├── service.go
│       └── taskService.go
├── static
│   └── upload.html
└── utils
    ├── aws
    │   └── s3.go
    ├── gorm
    │   └── gorm.go
    ├── log
    │   ├── log.go
    │   └── logByDate.go
    ├── paginator
    │   └── paginator.go
    ├── redis
    │   └── redis.go
    └── response
        └── response.go
```

# Folder definition
- Config
> Environment setting

- Controller
> Receiving HTTP requests calling requests and services

- Docker
> Project configurations

- Docs
> Swagger API documentation location

- Log
> Location of test log output files

- Middleware
> Intermediary layer, responsible for filtering incoming data

- Migration
> Create datatable details

- Model
> As a returned object

- Repository
> Assist service in calling sql query

- Request
> Assist controller validation request parameters

- Router
> API route locations

- Service
> Assist controller with business logic

- Static
> Test the API for the presence of CORS issues

- Utils
> Modular code placement for project calls