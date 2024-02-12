-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `permissions` (
                       `id` INT PRIMARY KEY AUTO_INCREMENT,
                       `title` VARCHAR(191) NOT NULL UNIQUE,
                       `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `permissions`;
