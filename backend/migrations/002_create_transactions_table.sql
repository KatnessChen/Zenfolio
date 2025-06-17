-- Migration: 002_create_transactions_table
-- Description: Create transactions table with comprehensive financial transaction tracking
-- Created: 2025-06-16

CREATE TABLE IF NOT EXISTS transactions (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
