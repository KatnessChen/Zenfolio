-- Useful SQL Queries for Transaction Tracker Database
-- Copy and paste these into Sequel Ace's Query tab or MySQL CLI
-- 

-- 1. Check all users
SELECT id, username, email, first_name, last_name, is_active, created_at, updated_at
FROM users
ORDER BY created_at;

-- 2. Check all transactions
SELECT id, user_id, type, symbol, quantity, price, amount, currency, 
       broker, transaction_date, user_notes, created_at
FROM transactions
ORDER BY transaction_date DESC;

-- 3. Get user transaction summary
SELECT 
    u.username,
    u.email,
    COUNT(t.id) as total_transactions,
    SUM(CASE WHEN t.type = 'buy' THEN t.amount ELSE 0 END) as total_buys,
    SUM(CASE WHEN t.type = 'sell' THEN t.amount ELSE 0 END) as total_sells,
    SUM(CASE WHEN t.type = 'dividend' THEN t.amount ELSE 0 END) as total_dividends
FROM users u
LEFT JOIN transactions t ON u.id = t.user_id
GROUP BY u.id, u.username, u.email;

-- 4. Portfolio holdings by user
SELECT 
    u.username,
    t.symbol,
    SUM(CASE WHEN t.type = 'buy' THEN t.quantity ELSE 0 END) as total_bought,
    SUM(CASE WHEN t.type = 'sell' THEN t.quantity ELSE 0 END) as total_sold,
    (SUM(CASE WHEN t.type = 'buy' THEN t.quantity ELSE 0 END) - 
     SUM(CASE WHEN t.type = 'sell' THEN t.quantity ELSE 0 END)) as current_holdings
FROM users u
JOIN transactions t ON u.id = t.user_id
WHERE t.type IN ('buy', 'sell')
GROUP BY u.id, u.username, t.symbol
HAVING current_holdings > 0
ORDER BY u.username, t.symbol;

-- 5. Recent transactions (last 30 days)
SELECT 
    u.username,
    t.type,
    t.symbol,
    t.quantity,
    t.price,
    t.amount,
    t.currency,
    t.broker,
    t.transaction_date,
    t.user_notes
FROM transactions t
JOIN users u ON t.user_id = u.id
WHERE t.transaction_date >= DATE_SUB(NOW(), INTERVAL 30 DAY)
ORDER BY t.transaction_date DESC;

-- 6. Check database schema
SHOW TABLES;

-- 7. Check table structure
DESCRIBE users;
DESCRIBE transactions;

-- 8. Check indexes
SHOW INDEX FROM users;
SHOW INDEX FROM transactions;

-- 9. Check migration status
SELECT * FROM schema_migrations ORDER BY version;

-- 10. Database statistics
SELECT 
    'users' as table_name,
    COUNT(*) as row_count
FROM users
UNION ALL
SELECT 
    'transactions' as table_name,
    COUNT(*) as row_count
FROM transactions;

-- 11. Transaction type distribution
SELECT 
    type,
    COUNT(*) as count,
    SUM(amount) as total_amount,
    AVG(amount) as avg_amount,
    ROUND((COUNT(*) * 100.0 / (SELECT COUNT(*) FROM transactions)), 2) as percentage
FROM transactions
GROUP BY type
ORDER BY count DESC;

-- 12. Top traded symbols
SELECT 
    symbol,
    COUNT(*) as transaction_count,
    SUM(CASE WHEN type = 'buy' THEN quantity ELSE 0 END) as total_bought,
    SUM(CASE WHEN type = 'sell' THEN quantity ELSE 0 END) as total_sold,
    SUM(amount) as total_value
FROM transactions
GROUP BY symbol
ORDER BY transaction_count DESC
LIMIT 10;

-- 13. Portfolio value by user (current holdings)
SELECT 
    u.username,
    t.symbol,
    (SUM(CASE WHEN t.type = 'buy' THEN t.quantity ELSE 0 END) - 
     SUM(CASE WHEN t.type = 'sell' THEN t.quantity ELSE 0 END)) as current_holdings,
    AVG(CASE WHEN t.type = 'buy' THEN t.price ELSE NULL END) as avg_buy_price,
    MAX(t.transaction_date) as last_transaction_date
FROM users u
JOIN transactions t ON u.id = t.user_id
WHERE t.type IN ('buy', 'sell')
GROUP BY u.id, u.username, t.symbol
HAVING current_holdings > 0
ORDER BY u.username, current_holdings DESC;

-- 14. Monthly transaction summary
SELECT 
    DATE_FORMAT(transaction_date, '%Y-%m') as month,
    COUNT(*) as transaction_count,
    SUM(CASE WHEN type = 'buy' THEN amount ELSE 0 END) as buy_volume,
    SUM(CASE WHEN type = 'sell' THEN amount ELSE 0 END) as sell_volume,
    SUM(CASE WHEN type = 'dividend' THEN amount ELSE 0 END) as dividend_income
FROM transactions
GROUP BY DATE_FORMAT(transaction_date, '%Y-%m')
ORDER BY month DESC;

-- 15. User activity summary
SELECT 
    u.username,
    u.email,
    COUNT(t.id) as total_transactions,
    MIN(t.transaction_date) as first_transaction,
    MAX(t.transaction_date) as last_transaction,
    SUM(t.amount) as total_transaction_value,
    COUNT(DISTINCT t.symbol) as unique_symbols_traded
FROM users u
LEFT JOIN transactions t ON u.id = t.user_id
GROUP BY u.id, u.username, u.email
ORDER BY total_transactions DESC;

-- 16. Broker usage analysis
SELECT 
    COALESCE(broker, 'Unknown') as broker,
    COUNT(*) as transaction_count,
    COUNT(DISTINCT user_id) as user_count,
    SUM(amount) as total_volume,
    AVG(amount) as avg_transaction_size
FROM transactions
GROUP BY broker
ORDER BY transaction_count DESC;

-- 17. Currency distribution
SELECT 
    currency,
    COUNT(*) as transaction_count,
    SUM(amount) as total_amount,
    ROUND((COUNT(*) * 100.0 / (SELECT COUNT(*) FROM transactions)), 2) as percentage
FROM transactions
GROUP BY currency
ORDER BY transaction_count DESC;

-- 18. Performance queries - Index usage check
-- Run EXPLAIN on common queries to verify index usage
EXPLAIN SELECT * FROM transactions WHERE user_id = 1;
EXPLAIN SELECT * FROM transactions WHERE user_id = 1 AND transaction_date > '2025-01-01';
EXPLAIN SELECT * FROM transactions WHERE symbol = 'AAPL' AND transaction_date > '2025-01-01';

-- 19. Data integrity checks
-- Check for orphaned transactions (shouldn't exist with foreign key)
SELECT t.id, t.user_id 
FROM transactions t 
LEFT JOIN users u ON t.user_id = u.id 
WHERE u.id IS NULL;

-- Check for invalid transaction types
SELECT DISTINCT type FROM transactions WHERE type NOT IN ('buy', 'sell', 'dividend');

-- Check for negative quantities on buy transactions
SELECT id, user_id, symbol, quantity, type 
FROM transactions 
WHERE type = 'buy' AND quantity <= 0;

-- 20. Cleanup queries (use with caution)
-- Remove soft-deleted records (if needed)
-- DELETE FROM transactions WHERE deleted_at IS NOT NULL;
-- DELETE FROM users WHERE deleted_at IS NOT NULL;
