CREATE TABLE IF NOT EXISTS orders(
    id char(36) PRIMARY KEY,
    customer_id char(36) NOT NULL,
    seller_id char(36) NOT NULL,
    costume_id int NOT NULL,
    shipment_destination varchar(30) NOT NULL,
    shipment_origin varchar(30) NOT NULL,
    total decimal(10,2) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE
);
