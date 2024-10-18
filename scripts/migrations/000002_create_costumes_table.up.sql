CREATE TABLE IF NOT EXISTS costumes(
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    name varchar(100) NOT NULL,
    description text NOT NULL,
    price decimal(10,2) NOT NULL,
    picture varchar(255),
    available bool DEFAULT TRUE,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
