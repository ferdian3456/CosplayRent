CREATE TABLE IF NOT EXISTS orders(
    id uuid PRIMARY KEY,
    user_id UUID NOT NULL,
    costume_id int NOT NULL,
    shipping_id int NOT NULL,
    total decimal(10,2) NOT NULL,
    status int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE,
    FOREIGN KEY (shipping_id) REFERENCES shippings(id) ON DELETE CASCADE
)
