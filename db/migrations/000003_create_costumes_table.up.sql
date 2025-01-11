CREATE TABLE IF NOT EXISTS costumes(
    id serial PRIMARY KEY,
    user_id char(36) NOT NULL,
    name varchar(100) NOT NULL,
    description text NOT NULL,
    material varchar(30) NOT NULL,
    size varchar(4) NOT NULL,
    weight int NOT NULL,
    category_id int NOT NULL,
    price decimal(10,2) NOT NULL,
    costume_picture varchar(255) NOT NULL,
    available VARCHAR(13) DEFAULT 'Ready',
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
