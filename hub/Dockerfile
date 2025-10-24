FROM golang:1.25-alpine

RUN go install github.com/air-verse/air@latest

RUN apk add --no-cache bash

WORKDIR /app

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
