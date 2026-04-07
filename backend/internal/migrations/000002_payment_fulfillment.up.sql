CREATE UNIQUE INDEX IF NOT EXISTS inventory_user_sku_idx
    ON inventory (user_id, sku);

CREATE TABLE IF NOT EXISTS payment_orders (
    external_id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'new',
    transaction_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS payment_order_items (
    external_id TEXT NOT NULL REFERENCES payment_orders(external_id) ON DELETE CASCADE,
    sku TEXT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (external_id, sku)
);

CREATE TABLE IF NOT EXISTS processed_payment_events (
    event_key TEXT PRIMARY KEY,
    external_id TEXT NOT NULL REFERENCES payment_orders(external_id) ON DELETE CASCADE,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
