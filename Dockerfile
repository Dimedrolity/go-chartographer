# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

CMD make ARGS="./data"

