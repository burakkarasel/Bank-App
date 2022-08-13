# _Cactus Bank_

![Project Image](./Cactus-Bank.gif)

---

### Table of Contents

- [Description](#description)
- [How To Use](#how-to-use)
- [Author Info](#author-info)

---

## Description

Cactus Bank enables it's customer to create user, login to those user. You can create multiple account for different currencies. You can deposit, withdraw to your account or transfer other accounts!

## Technologies

### Main Technologies

- [Go](https://go.dev/)
- [Gin Framework](https://github.com/gin-gonic/gin)
- [Github Actions](https://github.com/features/actions)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [gRPC](https://grpc.io/)

### Libraries

- [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [golang/mock](https://github.com/golang/mock)
- [golang/protobuf](https://github.com/golang/protobuf)
- [google/uuid](https://github.com/google/uuid)
- [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway/)
- [lib/pq](https://github.com/lib/pq)
- [o1egl/paseto](https://github.com/o1egl/paseto)
- [spf13/viper](https://github.com/spf13/viper)
- [stretchr/testify](https://github.com/stretchr/testify)
- [crypto](https://golang.org/x/crypto)
- [genproto](https://google.golang.org/genproto)
- [grpc](https://google.golang.org/grpc)
- [protoc-gen-go-grpc](https://google.golang.org/grpc/cmd/protoc-gen-go-grpc)
- [protobuf](https://google.golang.org/protobuf)

[Back To The Top](#cactus-bank)

---

## How To Use

### Tools

- [Go](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [TablePlus](https://tableplus.com/download)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)

### Setup Database

- Create bank-network

```
make network
```

- Start postgres container

```
make postgres
```

- Create bank_app Database

```
make createdb
```

- Migrations up

```
make up
```

- Migrations down

```
make down
```

### Generate Database functions

- Generate SQL CRUD functions

```
make sqlc
```

- Generate mockdb

```
make mock
```

### Run tests

```
make test
```

### Start App

- Start app directly

```
make server
```

- Run docker container

```
docker compose up
```

### Give it a try

#### Routes

| Request        | URL                                               | json                                                                       | Auth needed |
| -------------- | ------------------------------------------------- | -------------------------------------------------------------------------- | ----------- |
| Create user    | :8080/users                                       | {"username": "", "password": "", "email", "" "full_name": ""}              | No          |
| Login user     | :8080/users/login                                 | {"username": "", "password": ""}                                           | No          |
| Create account | :8080/accounts                                    | {"currency": ""}                                                           | Yes         |
| Get account    | :8080/accounts/:id                                |                                                                            | Yes         |
| List accounts  | :8080/accounts?page_id=1&page_size=5              |                                                                            | Yes         |
| Create entry   | :8080/entries                                     | {"account_id": 0, "amount":0}                                              | Yes         |
| Get entry      | :8080/entries/:id                                 |                                                                            | Yes         |
| List entries   | :8080/accounts?account_id=1&page_id=1&page_size=5 |                                                                            | Yes         |
| Make transfer  | :8080/transfers                                   | {"from_account_id": 0, "to_account_id": 0, "amount": 0, "currency": "USD"} | Yes         |

Don't forget to copy your access token for authentication required routes after logging in!

[Back To The Top](#cactus-bank)

---

## Author Info

- Twitter - [@dev_bck](https://twitter.com/dev_bck)

[Back To The Top](#cactus-bank)
