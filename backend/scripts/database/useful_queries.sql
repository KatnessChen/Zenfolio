-- Useful SQL Queries for Transaction Tracker Database
-- Copy and paste these into Sequel Ace's Query tab

-- 1. Check all users
SELECT id, username, email, first_name, last_name, is_active, created_at, last_login_at
FROM users
ORDER BY created_at;

-- 2. Check all transactions
SELECT id, user_id, type, symbol, quantity, price, amount, fee, currency, 
       broker, account, transaction_date, description, status
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
    t.transaction_date,
    t.description
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
