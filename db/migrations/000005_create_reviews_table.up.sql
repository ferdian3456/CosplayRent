CREATE TABLE IF NOT EXISTS reviews(
    id serial PRIMARY KEY,
    customer_id char(36) NOT NULL,
    costume_id int NOT NULL,
    order_id char(36) NOT NULL,
    description text NOT NULL,
    review_picture varchar(255) NOT NULL,
    rating int NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);
