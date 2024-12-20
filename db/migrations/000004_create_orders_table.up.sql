CREATE TABLE IF NOT EXISTS orders(
    id uuid PRIMARY KEY,
    user_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    costume_id int NOT NULL,
    total decimal(10,2) NOT NULL,
    status varchar(35) default 'Pending',
    description text,
    status_payment bool default false,
    status_shipping varchar(20) NOT NULL default false,
    is_cancelled bool default false,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE
);
