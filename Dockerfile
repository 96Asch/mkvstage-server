FROM golang:alpine as builder

LABEL version="1.0"
LABEL maintainer="Andrew Huang <aschhuang@gmail.com>"

RUN apk update && apk upgrade --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD [ "./server" ]
