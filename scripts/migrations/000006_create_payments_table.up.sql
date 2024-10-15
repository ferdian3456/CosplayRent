CREATE TABLE IF NOT EXISTS payments(
    id int PRIMARY KEY,
    order_id int NOT NULL,
    payment_method_id int NOT NULL,
    status int NOT NULL DEFAULT 0,
    total decimal(10,2) NOT NULL,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id) ON DELETE CASCADE
)
