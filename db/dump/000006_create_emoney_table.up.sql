CREATE TABLE IF NOT EXISTS emoney(
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    emoney_amount decimal(10,2) DEFAULT 0,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
