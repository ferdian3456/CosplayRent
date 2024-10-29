CREATE TABLE IF NOT EXISTS reviews(
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    costume_id int NOT NULL,
    description text NOT NULL,
    rating int NOT NULL DEFAULT NULL,
    created_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE
)
