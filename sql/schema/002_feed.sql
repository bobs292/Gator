-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_fetch_at TIMESTAMP,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
    );
-- +goose Down
DROP TABLE feeds;