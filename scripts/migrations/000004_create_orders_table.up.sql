CREATE TABLE IF NOT EXISTS orders(
    id uuid PRIMARY KEY,
    user_id UUID NOT NULL,
    costume_id int NOT NULL,
    shipping_id int default 0,
    total decimal(10,2) NOT NULL,
    status_payment bool default false,
    is_cancelled bool default false,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE
)
