# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /app

COPY . .

EXPOSE 8080

CMD make ARGS="./data"

