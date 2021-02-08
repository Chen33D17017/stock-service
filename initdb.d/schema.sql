DROP DATABASE IF EXISTS `stock_price`;
CREATE DATABASE `stock_price`;
USE `stock_price`;

SET NAMES utf8;
set character_set_client = utf8mb4;

CREATE TABLE `stock` (
    `id` INT AUTO_INCREMENT,
    `code` CHAR(50) UNIQUE NOT NULL,
    `name` CHAR(50),
    PRIMARY KEY (`id`)
)  ENGINE=INNODB;


CREATE TABLE`stock_data` (
	`stock_id` INT,
    `price_at` date,
    `open` double NOT NULL,
    `high` double NOT NULL,
    `low` double NOT NULL,
    `close` double NOT NULL,
    `vol` double NOT NULL,
    CONSTRAINT `contacts_pk` PRIMARY KEY (`stock_id`, `price_at`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
) ENGINE=INNODB;

CREATE TABLE `user` (
    `id` INT AUTO_INCREMENT,
    `username` CHAR(255) NOT NULL,
    `email` VARCHAR(255) UNIQUE NOT NULL,
    `password` binary(60) NOT NULL,
    PRIMARY KEY(`id`)
) ENGINE=INNODB;

CREATE TABLE `api`(
    `id` INT AUTO_INCREMENT,
    `type` CHAR(50)
    PRIMARY KEY(`id`)
) ENGINE=INNODB;

CREATE TABLE `user_api` (
    `id` INT AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `api_id` INT NOT NULL,
    `api_key` CHAR(255) NOT NULL,
    `api_secret` CHAR(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`api_id`) REFERENCES `api`(`id`)
) ENGINE=INNODB;

CREATE TABLE `stock_log`(
    `id` INT AUTO_INCREMENT,
    `stock_id` INT NOT NULL,
    `price` FLOAT NOT NULL,
    `time` TIMESTAMP NOT NULL
    PRIMARY KEY (`id`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
) ENGINE=INNODB;

CREATE TABLE `holiday` (
    `id` INT AUTO_INCREMENT,
    `name` CHAR(255),
    `date` DATE UNIQUE NOT NULL,
    PRIMARY KEY (`id`)
) 

-- cross_direction up: TRUE down: FALSE
CREATE TABLE `stock_alert` (
    `id` INT AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `stock_id` INT NOT NULL,
    `cross_direction` BOOLEAN NOT NULL,
    `price` FLOAT NOT NULL,
    `alert_on` BOOLEAN,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`),
    FOREIGN KEY (`stredegy_id`) REFERENCES `stredy`(`id`)
) ENGINE=INNODB;


-- buy_sell buy: TRUE sell: FALSE
CREATE TABLE `trade_log`(
    `id` INT AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `buy_sell` BOOLEAN NOT NULL,
    `time` TIMESTAMP NOT NULL,
    `price` FLOAT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
)

CREATE TABLE `cron_hook`(
    `id` INT AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `stock_id` INT NOT NULL,
    `price` FLOAT INT NOT NULL,
    PRIMARY KEY(`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
)