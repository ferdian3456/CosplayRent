CREATE TABLE IF NOT EXISTS payments(
    id serial PRIMARY KEY,
    order_id char(36) NOT NULL,
    customer_id char(36) NOT NULL,
    seller_id char(36) NOT NULL,
    status varchar(9) NOT NULL,
    amount decimal(10,2) NOT NULL,
    method varchar(15) NOT NULL,
    midtrans_redirect_url varchar(255),
    midtrans_url_expired_time timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);