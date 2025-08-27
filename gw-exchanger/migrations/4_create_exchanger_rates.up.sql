CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate DECIMAL(10, 4) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_currency_pair UNIQUE (from_currency, to_currency)
);

INSERT INTO exchange_rates (from_currency, to_currency, rate)
VALUES
    ('USD', 'RUB', 75.50),
    ('USD', 'EUR', 0.85),
    ('EUR', 'RUB', 89.25)
ON CONFLICT (from_currency, to_currency) DO NOTHING;
