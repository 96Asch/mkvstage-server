FROM golang:alpine as testBuilder

LABEL version="1.0"
LABEL maintainer="Andrew Huang <aschhuang@gmail.com>"

RUN apk update && apk upgrade --no-cache git

WORKDIR /test

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN command


CMD CGO_ENABLED=0 go test ./...  -coverpkg=./... -coverprofile ./coverage.out && go tool cover -func ./coverage.out
