-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS kyc_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(20, 8) NOT NULL,
    mpesa_number VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    status ENUM('pending', 'confirmed', 'failed') DEFAULT 'pending',
    mpesa_transaction_code VARCHAR(100),
    mpesa_checkout_id VARCHAR(50),
    mpesa_receipt_number VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kyc_payments;
-- +goose StatementEnd