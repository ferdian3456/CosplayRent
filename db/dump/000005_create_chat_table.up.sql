CREATE TABLE IF NOT EXISTS chats(
    id serial PRIMARY KEY,
    sender uuid NOT NULL,
    receiver uuid NOT NULL,
    message text NOT NULL,
    created_at timestamp NOT NULL
)
