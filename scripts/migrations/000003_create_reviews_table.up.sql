CREATE TABLE IF NOT EXISTS reviews(
    id int PRIMARY KEY,
    user_id uuid NOT NULL,
    costume_id int NOT NULL,
    rating int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE
)
