CREATE TABLE IF NOT EXISTS orders(
    id uuid PRIMARY KEY,
    user_id UUID NOT NULL,
    costume_id int NOT NULL,
    total decimal(10,2) NOT NULL,
    status_payment bool default false,
    status_shipping varchar(20) NOT NULL,
    is_cancelled bool default false,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE
)
