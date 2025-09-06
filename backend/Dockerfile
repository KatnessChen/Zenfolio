# syntax=docker/dockerfile:1
FROM golang:1.21-alpine
WORKDIR /app
RUN apk add --no-cache git tzdata
RUN go install github.com/cosmtrek/air@v1.41.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.2
COPY . .
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
