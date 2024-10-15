CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY,
    name varchar(20) NOT NULL,
    email varchar(30),
    address varchar(100),
    password varchar(20) NOT NULL,
    role int NOT NULL,
    profile_picture varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL DEFAULT NULL
);
