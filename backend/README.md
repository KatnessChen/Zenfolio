# Transaction Tracker API

A RESTful Golang API for tracking financial transactions.

## Features

- RESTful API with Gin framework
- JWT Authentication
- Rate limiting
- OpenAPI documentation
- Example GET endpoint for testing

## Prerequisites

- Go 1.16 or higher
- Git

## Getting Started

### Clone the repository

```bash
git clone https://github.com/your-username/transaction-tracker.git
cd transaction-tracker/backend
```

### Install dependencies

```bash
go mod download
```

### Set environment variables (optional)

```bash
export SERVER_ADDR=":8080"
export JWT_SECRET="your-secret-key"
```

### Run the application

```bash
go run main.go
```

The server will start on http://localhost:8080.

## API Documentation

OpenAPI documentation is available in `/docs/api.yaml`.

### Example Endpoints

#### Health Check

```
GET /health
```

#### Authentication

```
POST /api/v1/login
{
  "username": "user123",
  "password": "password123"
}
```

#### Hello World (Protected)

```
GET /api/v1/hello-world
Authorization: Bearer <jwt_token>
```

## Rate Limiting

The API implements rate limiting to prevent abuse. By default, it allows 100 requests per minute per user.

## Authentication

JWT-based authentication is implemented. To access protected endpoints:

1. Get a token using the `/api/v1/login` endpoint
2. Include the token in the Authorization header: `Authorization: Bearer <token>`
