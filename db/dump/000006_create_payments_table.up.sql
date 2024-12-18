CREATE TABLE IF NOT EXISTS payments(
    id serial PRIMARY KEY,
    payment_method_id int NOT NULL,
    status int NOT NULL DEFAULT 0,
    total decimal(10,2) NOT NULL,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id) ON DELETE CASCADE
)
