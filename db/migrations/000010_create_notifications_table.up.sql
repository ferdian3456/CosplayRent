CREATE TABLE IF NOT EXISTS notifications(
    id serial PRIMARY KEY,
    template_name varchar(255) NOT NULL,
    template_subject varchar(255) NOT NULL,
    template_body text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
)