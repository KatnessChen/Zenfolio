# Transaction Tracker Database Setup Guide

## Overview

This comprehensive guide covers the complete database setup and management for the Transaction Tracker application. The database layer provides a robust, scalable foundation for data persistence using MySQL and GORM.

## Table of Contents

1. [Database Architecture](#database-architecture)
2. [Database Configuration](#database-configuration)
3. [Database Setup Methods](#database-setup-methods)
   - [Automated Setup (Using Scripts)](#automated-setup-using-scripts)
   - [Manual Setup](#manual-setup)
4. [Database CLI Tool Usage](#database-cli-tool-usage)
5. [Testing and Verification](#testing-and-verification)
6. [Advanced Usage](#advanced-usage)
7. [Troubleshooting](#troubleshooting)

---

## Database Architecture

### Components

1. **Database Configuration** (`config/database.go`)

   - Environment-specific database settings
   - Connection string management
   - SSL/TLS configuration

2. **Database Models** (`internal/models/models.go`)

   - User model for authentication and profiles
   - Transaction model for financial transactions
   - Base model with common fields (ID, timestamps, soft delete)

3. **Database Connection Manager** (`internal/database/database.go`)

   - Connection establishment and pooling
   - Health check mechanisms
   - Graceful shutdown procedures
   - Retry logic with exponential backoff

4. **Migration System** (`migrations/`)

   - Schema versioning and management
   - Migration runner with rollback capabilities
   - Automated migration tracking

5. **Database Services** (`internal/services/`)

   - User service for user-related operations
   - Transaction service for transaction management
   - Repository pattern implementation

6. **Database Seeder** (`internal/database/seeder.go`)
   - Development data seeding
   - Test data management
   - Environment-specific seed data

### Database Schema

#### Users Table

```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_active (is_active),
    INDEX idx_users_deleted_at (deleted_at)
);
```

#### Transactions Table

```sql
CREATE TABLE transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    type VARCHAR(50) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    quantity DECIMAL(15,4) NOT NULL,
    price DECIMAL(15,4) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    fee DECIMAL(15,2) DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    broker VARCHAR(100),
    account VARCHAR(100),
    transaction_date TIMESTAMP NOT NULL,
    settlement_date TIMESTAMP NULL,
    description TEXT,
    reference VARCHAR(255),
    status VARCHAR(20) DEFAULT 'completed',
    tags TEXT,
    metadata JSON,
    extracted_from VARCHAR(255),
    processing_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_transactions_user_id (user_id),
    INDEX idx_transactions_type (type),
    INDEX idx_transactions_symbol (symbol),
    INDEX idx_transactions_date (transaction_date),
    INDEX idx_transactions_status (status),
    INDEX idx_transactions_user_date (user_id, transaction_date),
    INDEX idx_transactions_symbol_date (symbol, transaction_date),
    INDEX idx_transactions_type_status (type, status),
    INDEX idx_transactions_deleted_at (deleted_at),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT
);
```

---

## Database Configuration

### Environment Variables

Before setting up the database, you need to configure environment variables in your `.env` file:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_NAME=transaction_tracker_dev
DB_USER=tracker_user
DB_PASSWORD=tracker_password
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE=10
DB_CONN_MAX_LIFETIME=3600
DB_CONN_MAX_IDLE_TIME=1800
DB_CHARSET=utf8mb4
DB_LOC=Local

# Server Configuration
SERVER_ADDR=:8080

# JWT Configuration
JWT_SECRET=dev-jwt-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24

# Rate Limiting
RATE_LIMIT_REQUESTS=100

# AI Model Configuration
AI_MODEL=gemini-1.5-flash
GEMINI_API_KEY=your-gemini-api-key-here
AI_TIMEOUT=30
```

### Multi-Environment Support

The database layer supports multiple environments:

- **Development**: Local MySQL with minimal security (`transaction_tracker_dev`)
- **Test**: Isolated test database (`transaction_tracker_test`)
- **Staging**: Remote MySQL with SSL required
- **Production**: Clustered MySQL with full security

---

## Database Setup Methods

### Prerequisites

Before starting, ensure you have:

- [ ] Go 1.23.1+ installed
- [ ] Git installed
- [ ] MySQL 8.0+ OR Docker installed
- [ ] Terminal/Command line access
- [ ] Text editor for configuration

### Automated Setup (Using Scripts)

#### Quick Start with Setup Script

The fastest way to get started is using the automated setup script:

```bash
# Navigate to backend directory
cd /Users/kc/Develop/100_Case/04_Transaction_Tracker/backend

# Make script executable
chmod +x scripts/database/setup-database.sh

# Run automated setup
./scripts/database/setup-database.sh
```

**What the script does:**

1. Checks for Docker availability
2. Starts MySQL container with docker-compose or docker run
3. Waits for MySQL to be ready
4. Tests database connection
5. Runs migrations
6. Seeds sample data

#### Using Docker Compose (Recommended)

```bash
# Start MySQL container
docker-compose up -d mysql

# Check if container is running
docker-compose ps

# Check MySQL logs (wait for "ready for connections")
docker-compose logs mysql
```

### Manual Setup

#### Option 1: Native MySQL Installation

1. **Install MySQL**:

   ```bash
   # macOS with Homebrew
   brew install mysql
   brew services start mysql

   # Set root password
   mysql_secure_installation
   ```

2. **Create Database and User**:

   ```bash
   # Login to MySQL as root
   mysql -u root -p
   ```

   ```sql
   -- Create databases
   CREATE DATABASE transaction_tracker_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE DATABASE transaction_tracker_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

   -- Create user
   CREATE USER 'tracker_user'@'localhost' IDENTIFIED BY 'tracker_password';

   -- Grant privileges
   GRANT ALL PRIVILEGES ON transaction_tracker_dev.* TO 'tracker_user'@'localhost';
   GRANT ALL PRIVILEGES ON transaction_tracker_test.* TO 'tracker_user'@'localhost';

   -- Flush privileges
   FLUSH PRIVILEGES;

   -- Exit
   EXIT;
   ```

#### Option 2: Docker Manual Setup

```bash
# Run MySQL container manually
docker run -d \
  --name transaction_tracker_mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=rootpassword \
  -e MYSQL_DATABASE=transaction_tracker_dev \
  -e MYSQL_USER=tracker_user \
  -e MYSQL_PASSWORD=tracker_password \
  mysql:8.0 \
  --default-authentication-plugin=mysql_native_password

# Wait for MySQL to start (about 30 seconds)
docker logs transaction_tracker_mysql

# Test connection
docker exec -it transaction_tracker_mysql mysql -u tracker_user -p transaction_tracker_dev
```

#### Configuration Setup

1. **Create Environment File**:

   ```bash
   # Copy the example environment file
   cp .env.example .env
   ```

2. **Edit Configuration** (update `.env` with your settings)

---

## Database CLI Tool Usage

The `db-cli.go` tool provides comprehensive database management capabilities.

### Basic Usage

```bash
# General syntax
go run tools/db-cli.go -action=<action> [-env=<environment>]

# Or build and use
go build -o db-cli tools/db-cli.go
./db-cli -action=<action>
```

### Available Actions

#### Health Check

```bash
# Basic health check
go run tools/db-cli.go -action=health
```

**Expected Output:**

```
✅ Database connection is healthy
📊 Connection Stats:
   Max Open Connections: 100
   Open Connections: 1
   In Use: 1
   Idle: 0
```

#### Migration Management

```bash
# Run all pending migrations
go run tools/db-cli.go -action=migrate

# Check migration status
go run tools/db-cli.go -action=status

# Rollback last migration
go run tools/db-cli.go -action=rollback
```

**Migration Status Output:**

```
Migration Status:
- 001_create_users_table.sql: ✅ Applied
- 002_create_transactions_table.sql: ✅ Applied
```

#### Data Seeding

```bash
# Seed development data (default)
go run tools/db-cli.go -action=seed

# Seed for specific environment
go run tools/db-cli.go -action=seed -env=development
go run tools/db-cli.go -action=seed -env=staging
go run tools/db-cli.go -action=seed -env=production  # Requires confirmation
```

**Seeding Output:**

```
🌱 Seeding database for environment: development
✅ Database seeded successfully!

Created:
- 3 sample users
- 15 sample transactions
- Portfolio data initialized
```

### CLI Tool Environment Variables

The CLI tool uses the same environment variables as the main application. Ensure your `.env` file is properly configured before running any CLI commands.

---

## Testing and Verification

### Step-by-Step Verification

#### 1. Test Database Connection

```bash
go run tools/db-cli.go -action=health
```

#### 2. Run Migrations

```bash
go run tools/db-cli.go -action=migrate
go run tools/db-cli.go -action=status
```

#### 3. Seed Sample Data

```bash
go run tools/db-cli.go -action=seed
```

#### 4. Start Application

```bash
go run main.go
```

#### 5. Test API Endpoints

```bash
# Basic health check
curl http://localhost:8080/health

# Database health check
curl http://localhost:8080/api/v1/health/database

# Register user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login (save the JWT token)
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 6. Test Transaction Operations

```bash
# Set JWT token from login response
export JWT_TOKEN="YOUR_JWT_TOKEN_HERE"

# Create transaction
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "type": "buy",
    "symbol": "AAPL",
    "quantity": 10,
    "price": 150.00,
    "date": "2025-06-17T10:00:00Z"
  }'

# Get user transactions
curl -H "Authorization: Bearer $JWT_TOKEN" \
  http://localhost:8080/api/v1/transactions

# Get portfolio
curl -H "Authorization: Bearer $JWT_TOKEN" \
  http://localhost:8080/api/v1/transactions/portfolio
```

### Automated Testing Script

Use the comprehensive testing script:

```bash
# Make script executable
chmod +x scripts/database/test-db-implementation.sh

# Run all tests
./scripts/database/test-db-implementation.sh
```

---

## Advanced Usage

### Using Services in Code

```go
// Initialize database
dm, err := database.Initialize(cfg)
if err != nil {
    log.Fatal(err)
}

// Create services
userService := services.NewUserService(dm.GetDB())
transactionService := services.NewTransactionService(dm.GetDB())

// Use user service
user := &models.User{
    Username: "john_doe",
    Email:    "john@example.com",
    // ... other fields
}
if err := userService.CreateUser(user); err != nil {
    log.Fatal(err)
}

// Use transaction service
transaction := &models.Transaction{
    UserID:   user.ID,
    Type:     "buy",
    Symbol:   "AAPL",
    Quantity: 100,
    Price:    150.25,
    Amount:   15025.00,
    // ... other fields
}
if err := transactionService.CreateTransaction(transaction); err != nil {
    log.Fatal(err)
}
```

### Health Check Endpoints

- `GET /health` - Basic application health
- `GET /api/v1/health/database` - Basic database health check
- `GET /api/v1/health/database/detailed` - Detailed database health with connection stats (requires authentication)

### Performance Monitoring

```bash
# Load testing with hey
go install github.com/rakyll/hey@latest

# Test health endpoint
hey -n 1000 -c 10 http://localhost:8080/health

# Test database health
hey -n 100 -c 5 http://localhost:8080/api/v1/health/database
```

### Security Features

- **Connection Security**: SSL/TLS encryption, environment variable-based credentials
- **Data Security**: Row-level security, soft delete, parameterized queries, bcrypt password hashing
- **Access Control**: Database user isolation, restricted foreign keys, audit trails

### Performance Optimizations

- **Indexing Strategy**: Primary keys, composite indexes, foreign key indexes
- **Connection Pooling**: Configurable settings, lifetime management, idle timeouts
- **Query Optimization**: Efficient JOINs, batch operations, pagination support

---

## Troubleshooting

### Common Issues and Solutions

#### 1. Connection Refused Error

```bash
# Check if MySQL is running
# For Docker:
docker-compose ps
docker-compose logs mysql

# For native MySQL:
brew services list | grep mysql
sudo systemctl status mysql
```

**Solution**: Start MySQL service or container.

#### 2. Access Denied Error

```bash
# Verify credentials
mysql -u tracker_user -p -h localhost transaction_tracker_dev
```

**Solution**: Check username/password in `.env` file or recreate user.

#### 3. Database Does Not Exist Error

```sql
-- Connect as root and create database
mysql -u root -p
CREATE DATABASE transaction_tracker_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 4. Migration Failed Error

```bash
# Check migration status
go run tools/db-cli.go -action=status

# Check for conflicts in migration files
```

#### 5. Port Already in Use

```bash
# Check what's using port 3306
lsof -i :3306

# Stop conflicting service
brew services stop mysql
# or
docker-compose down
```

### Database Console Access

```bash
# Via Docker
docker exec -it transaction_tracker_mysql mysql -u tracker_user -p transaction_tracker_dev

# Via native MySQL
mysql -u tracker_user -p transaction_tracker_dev
```

### Useful SQL Commands

```sql
-- Check tables
SHOW TABLES;

-- Check table structures
DESCRIBE users;
DESCRIBE transactions;

-- Count records
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM transactions;

-- Check migrations
SELECT * FROM schema_migrations;

-- Sample data queries
SELECT u.username, COUNT(t.id) as transaction_count
FROM users u
LEFT JOIN transactions t ON u.id = t.user_id
GROUP BY u.id;
```

---

## Success Criteria

You've successfully set up the database when:

- [ ] MySQL is running and accessible
- [ ] Database connection test passes (`go run tools/db-cli.go -action=health`)
- [ ] All migrations run successfully (`go run tools/db-cli.go -action=migrate`)
- [ ] Sample data is seeded (`go run tools/db-cli.go -action=seed`)
- [ ] Application starts without errors (`go run main.go`)
- [ ] Health check endpoints respond correctly
- [ ] User registration/login works
- [ ] Transaction CRUD operations work
- [ ] Database CLI commands execute successfully

---

## Additional Resources

- **API Documentation**: `docs/api.yaml`
- **Postman Collection**: `docs/postman_collection.json`
- **Migration Files**: `migrations/` directory
- **Test Suite**: `test/database_test.go`
- **Useful Queries**: `scripts/database/useful_queries.sql`

---

**🚀 Your database layer is now ready for development and testing!**

For additional support, check the troubleshooting section or review the implementation files in the repository.
