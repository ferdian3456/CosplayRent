CREATE TABLE IF NOT EXISTS user_verifications(
    id serial PRIMARY KEY,
    user_id char(36) NOT NULL,
    verification_code char(5) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)