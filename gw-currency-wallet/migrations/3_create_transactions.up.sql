CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    from_currency VARCHAR(3) NOT NULL,
    to_currency   VARCHAR(3),
    amount BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
