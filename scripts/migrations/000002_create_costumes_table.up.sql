CREATE TABLE IF NOT EXISTS costumes(
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    name varchar(100) NOT NULL,
    description text NOT NULL,
    bahan varchar(30) NOT NULL,
    ukuran varchar(30) NOT NULL,
    berat varchar(30) NOT NULL,
    kategori varchar(30) NOT NULL,
    price decimal(10,2) NOT NULL,
    costume_picture varchar(255) NOT NULL,
    available bool DEFAULT TRUE,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
