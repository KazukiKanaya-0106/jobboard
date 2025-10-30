FROM golang:1.25-alpine

RUN go install github.com/air-verse/air@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN apk add --no-cache bash

WORKDIR /app

EXPOSE 8080

COPY . .

RUN sqlc generate

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]
