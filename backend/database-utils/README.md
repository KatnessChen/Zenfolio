# Database Utilities

This folder contains database setup and utility scripts for the Transaction Tracker project.

## Files

### 🚀 Setup Scripts

#### `setup-database.sh`

**Purpose**: Automated database setup script for development environment

- Detects and starts MySQL with Docker or docker-compose
- Creates `.env` file if missing
- Runs database migrations and seeding
- Tests database connectivity

**Usage**:

```bash
cd /Users/kc/Develop/100_Case/04_Transaction_Tracker/backend
./database-utils/setup-database.sh
```

#### `test-db-implementation.sh`

**Purpose**: Comprehensive database implementation testing script

- Tests database connectivity
- Validates all migrations
- Runs seeding operations
- Performs CRUD operation tests

**Usage**:

```bash
./database-utils/test-db-implementation.sh
```

### 📊 Database Queries

#### `useful_queries.sql`

**Purpose**: Collection of useful SQL queries for database exploration

### 🔧 Database CLI Tool (Go)

#### `../tools/db-cli.go`

**Purpose**: Go-based database CLI for migrations, seeding, and health checks

**Usage**:

```bash
# Direct usage
go run tools/db-cli.go -action=health
go run tools/db-cli.go -action=migrate
go run tools/db-cli.go -action=seed

# Or build and use
go build -o db-cli tools/db-cli.go
./db-cli -action=status
```

**Available Actions**:

- `health` - Check database connectivity and stats
- `migrate` - Run pending migrations
- `rollback` - Rollback last migration
- `status` - Show migration status
- `seed` - Seed database with sample data

**Categories**:

- Basic data queries (users, transactions)
- Analytics queries (summaries, portfolios)
- Administrative queries (schema, indexes, migrations)

## Quick Start

1. **Setup Database**: `./database-utils/setup-database.sh`
2. **Test Implementation**: `./database-utils/test-db-implementation.sh`
3. **Explore Data**: Use queries from `useful_queries.sql` in Sequel Ace

## Connection Details

- Host: localhost:3306
- Database: transaction_tracker_dev
- User: tracker_user
- Password: tracker_password

## Troubleshooting

### Make scripts executable

```bash
chmod +x database-utils/*.sh
```

### Docker not running

```bash
open -a Docker
```

### Database connection issues

- Check container status: `docker-compose ps`
- Verify `.env` file exists with correct credentials
