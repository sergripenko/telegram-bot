
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users_locations(
    id          SERIAL PRIMARY KEY,
    user_id     INT,
    latitude    INT,
    longitude   INT,
    city        VARCHAR,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS users_locations;
