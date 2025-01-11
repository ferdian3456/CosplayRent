CREATE TABLE IF NOT EXISTS wishlists(
    id serial PRIMARY KEY,
    user_id char(36) NOT NULL,
    costume_id int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (costume_id) REFERENCES costumes(id) ON DELETE CASCADE
)