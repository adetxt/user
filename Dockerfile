FROM golang:alpine AS builder

WORKDIR /app

COPY . /app

RUN go get -u github.com/adetxt/protoc-gen-go-edison/cmd/protoc-gen-go-edison@v0.0.1
RUN go get github.com/bufbuild/buf/cmd/buf@v1.9.0

RUN go install github.com/adetxt/protoc-gen-go-edison/cmd/protoc-gen-go-edison@v0.0.1
RUN go install github.com/bufbuild/buf/cmd/buf@v1.9.0

RUN buf generate proto

ENV MYSQL_DB_HOST=103.150.196.254

RUN go build -o main . 

CMD ["/app/main"]
