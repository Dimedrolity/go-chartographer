# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS builder

WORKDIR /app

COPY . .

RUN apk add make && make clean build

# multistage-build для уменьшения итогового образа
FROM alpine

COPY --from=builder /app/build/app /

EXPOSE 8080

CMD [ "/app", "./data" ]