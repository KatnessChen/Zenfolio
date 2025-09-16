# syntax=docker/dockerfile:1
FROM golang:1.23-alpine
WORKDIR /app
RUN apk add --no-cache git tzdata
RUN apk add --no-cache gcc musl-dev
RUN go install github.com/cosmtrek/air@v1.41.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.2
RUN wget -O - https://github.com/golangci/golangci-lint/releases/download/v1.61.0/golangci-lint-1.61.0-linux-amd64.tar.gz | tar -xz -C /tmp \
	&& mv /tmp/golangci-lint-1.61.0-linux-amd64/golangci-lint /usr/local/bin/ \
	&& chmod +x /usr/local/bin/golangci-lint
COPY . .
EXPOSE 8081
CMD ["air", "-c", ".air.toml"]
