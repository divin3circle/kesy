-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS mints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    amount_kes BIGINT NOT NULL,
    wallet_address VARCHAR(255) NOT NULL,
    status ENUM('pending', 'completed', 'failed', 'settled') DEFAULT 'pending',
    restriction_end_date TIMESTAMPTZ,
    tokens_minted BIGINT,
    settled_txn_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
)
-- goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mints;
-- +goose StatementEnd