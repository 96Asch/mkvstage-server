FROM golang:alpine as base

RUN apk update && apk upgrade
RUN apk add curl


# Run the air command in the directory where our code will live
WORKDIR /opt/app/api

COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . /opt/app/api/

# Create another stage called "dev" that is based off of our "base" stage (so we have golang available to us)
FROM base as dev

# Install the air binary so we get live code-reloading when we save files
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin





CMD ["air"]

