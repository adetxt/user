# User

This app is using gRPC and Edison. Edison is my personal Go package that will simplify build gRPC server. Edison has ability to automate init gRPC gateway, so you can easily call the RPC method with gRPC and REST. All is handled by Edison.

## Getting started
- Download protoc-gen-go-edison `go install github.com/adetxt/protoc-gen-go-edison/cmd/protoc-gen-go-edison`
- Download buf `bufbuild/buf/buf` [toturial](https://docs.buf.build/installation)
- Generate code from proto `buf generate proto`
- `go mod tidy`
- `go run main.go`

## Structure

```
.
├── Dockerfile
├── README.md
├── buf.gen.yaml
├── config
│   └── config.go
├── domain
│   ├── auth.go
│   ├── general.go
│   └── user.go
├── gen
├── go.mod
├── go.sum
├── handler
│   └── grpc
│       ├── account_handler.go
│       └── auth_handler.go
├── main.go
├── proto
│   ├── account
│   │   └── v1
│   │       ├── account.proto
│   │       └── auth.proto
│   ├── buf.lock
│   └── buf.yaml
├── repository
│   └── user_mysql
│       ├── dto.go
│       └── user_mysql_repository.go
├── usecase
│   ├── auth_usecase.go
│   └── user_usecase.go
└── utils
    ├── auth
    │   └── jwt.go
    ├── mysql
    │   └── mysql.go
    └── password
        └── password.go
```

- **config** -- setup ENVAR config
- **domain** -- setup entity, and business process (usecase and repository interface)
- **handler** -- it's something like controller in MVC
- **proto** -- gRPC protobuf schemas
- **repository** -- repository implementation
- **usecase** -- usecase implementation
- **utils** -- helpers

## Credentials
```
Admin
email: admin@mail.com
password: password

User
email: user@mail.com
password: password
```

## API DOC
[https://documenter.getpostman.com/view/18749474/2s8YsxuWPK](https://documenter.getpostman.com/view/18749474/2s8YsxuWPK)