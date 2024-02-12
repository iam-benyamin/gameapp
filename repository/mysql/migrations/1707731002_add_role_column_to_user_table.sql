-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `users` ADD COLUMN `role` ENUM('user', 'admin') NOT NULL;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `role`;
