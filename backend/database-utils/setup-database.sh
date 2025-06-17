#!/bin/bash

# Database Setup Script for Transaction Tracker
# This script helps set up the MySQL database for development

set -e

# Change to backend directory (parent of database-utils)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
cd "$BACKEND_DIR"

echo "🚀 Transaction Tracker Database Setup"
echo "================================="
echo "Working directory: $(pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
  echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
  echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
  echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is available
if command -v docker &> /dev/null; then
  print_status "Docker found. Setting up MySQL with Docker..."

  # Check if docker-compose is available
  if command -v docker-compose &> /dev/null; then
    print_status "Starting MySQL container with docker-compose..."
    docker-compose up -d mysql

    print_status "Waiting for MySQL to be ready..."
    sleep 10

    # Wait for MySQL to be healthy
    while ! docker-compose exec mysql mysqladmin ping -h localhost --silent; do
      print_status "Waiting for MySQL to start..."
      sleep 2
    done

    print_status "MySQL is ready!"

  else
    print_status "Starting MySQL container with docker run..."
    docker run -d \
      --name transaction_tracker_mysql \
      -p 3306:3306 \
      -e MYSQL_ROOT_PASSWORD=rootpassword \
      -e MYSQL_DATABASE=transaction_tracker_dev \
      -e MYSQL_USER=tracker_user \
      -e MYSQL_PASSWORD=tracker_password \
      mysql:8.0 \
      --default-authentication-plugin=mysql_native_password

    print_status "Waiting for MySQL to be ready..."
    sleep 15

    print_status "MySQL container started!"
  fi

else
  print_warning "Docker not found. Please install MySQL manually:"
  echo ""
  echo "macOS (with Homebrew):"
  echo "  brew install mysql"
  echo "  brew services start mysql"
  echo ""
  echo "Ubuntu/Debian:"
  echo "  sudo apt update"
  echo "  sudo apt install mysql-server"
  echo "  sudo systemctl start mysql"
  echo ""
  echo "CentOS/RHEL:"
  echo "  sudo yum install mysql-server"
  echo "  sudo systemctl start mysqld"
  echo ""
  exit 1
fi

# Test database connection
print_status "Testing database connection..."

if [ -f ".env" ]; then
  print_status "Found .env file. Testing connection with Go application..."

  # Build and test the database CLI tool
  if [ -f "tools/db-cli.go" ]; then
    print_status "Building database CLI tool..."
    go build -o db-cli tools/db-cli.go

    # Export database credentials for the CLI tool
    export DB_USER=tracker_user
    export DB_PASSWORD=tracker_password

    print_status "Testing database connection..."
    ./db-cli -action=health

    print_status "Running database migrations..."
    ./db-cli -action=migrate

    print_status "Seeding database with sample data..."
    ./db-cli -action=seed

    # Clean up
    rm -f db-cli

    print_status "✅ Database setup complete!"
    echo ""
    echo "Your database is ready at:"
    echo "  Host: localhost"
    echo "  Port: 3306"
    echo "  Database: transaction_tracker_dev"
    echo "  User: tracker_user"
    echo ""
    echo "You can now start the application with:"
    echo "  go run main.go"

  else
    print_warning "Database CLI tool not found. Please run migrations manually."
  fi
else
  print_error ".env file not found. Please copy .env.example to .env and configure database settings."
  exit 1
fi
