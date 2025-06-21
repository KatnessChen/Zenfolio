-- Add deleted_at column to jwt_tokens table for soft delete support
ALTER TABLE jwt_tokens ADD COLUMN deleted_at TIMESTAMP NULL;

-- Add index for soft delete queries
CREATE INDEX idx_jwt_tokens_deleted_at ON jwt_tokens(deleted_at);
