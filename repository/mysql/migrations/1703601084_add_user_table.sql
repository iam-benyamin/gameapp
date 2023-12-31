-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users (
    id int primary key AUTO_INCREMENT,
    name varchar(255) not null ,
    phone_number varchar(255) not null unique,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE users;
