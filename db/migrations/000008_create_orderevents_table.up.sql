CREATE TABLE IF NOT EXISTS order_events(
    id serial PRIMARY KEY,
    user_id char(36) NOT NULL,
    order_id char(36) NOT NULL,
    status varchar(35) NOT NULL,
    notes text,
    shipment_receipt_user_id varchar(25) DEFAULT '',
    created_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);