
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users(
    id          SERIAL PRIMARY KEY,
    user_name   VARCHAR,
    first_name  VARCHAR,
    last_name   VARCHAR,
    chat_id     INT,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX chat_id ON users(chat_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS users;
DROP INDEX chat_id;
