CREATE TABLE IF NOT EXISTS costumes(
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    name varchar(100) NOT NULL,
    description text NOT NULL,
    price decimal(10,2) NOT NULL,
    picture varchar(255) NOT NULL,
    available bool NOT NULL DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
