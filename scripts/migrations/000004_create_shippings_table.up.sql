CREATE TABLE IF NOT EXISTS shippings(
    id int PRIMARY KEY,
    name varchar(30) NOT NULL,
    price decimal(10,2) NOT NULL,
    origin varchar(255) NOT NULL,
    destinantion varchar(255) NOT NULL
)
