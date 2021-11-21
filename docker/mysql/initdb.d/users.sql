CREATE DATABASE IF NOT EXISTS dev;
USE dev;

CREATE TABLE IF NOT EXISTS todos
(
    `id` INT NOT NULL AUTO_INCREMENT,
    `todo` VARCHAR(255) NOT NULL,
    `completed` TINYINT DEFAULT 0,
    PRIMARY KEY (`id`)
);

INSERT INTO todos
    (todo)
VALUES
    ('Golang'),
    ('React'),
    ('PHP'),
    ('Rust');
