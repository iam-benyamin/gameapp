-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users ADD COLUMN password VARCHAR(255) NOT NULL;

-- +migrate Down
ALTER TABLE users DROP COLUMN password;