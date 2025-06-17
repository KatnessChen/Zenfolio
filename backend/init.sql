-- Initialize Transaction Tracker Database
CREATE DATABASE IF NOT EXISTS transaction_tracker_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS transaction_tracker_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create user if not exists
CREATE USER IF NOT EXISTS 'tracker_user'@'%' IDENTIFIED BY 'tracker_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON transaction_tracker_dev.* TO 'tracker_user'@'%';
GRANT ALL PRIVILEGES ON transaction_tracker_test.* TO 'tracker_user'@'%';

-- Flush privileges
FLUSH PRIVILEGES;
