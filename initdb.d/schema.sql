DROP DATABASE IF EXISTS `stock_price`;
CREATE DATABASE `stock_price`;
USE `stock_price`;

SET NAMES utf8;
set character_set_client = utf8mb4;

CREATE TABLE `stock` (
    `id` CHAR(50) UNIQUE NOT NULL,
    `name` CHAR(50),
    PRIMARY KEY (`id`)
)  ENGINE=INNODB;


CREATE TABLE`stock_data` (
    `id` INT AUTO_INCREMENT,
	`stock_id` CHAR(50) NOT NULL,
    `price_at` date,
    `open` double NOT NULL,
    `high` double NOT NULL,
    `low` double NOT NULL,
    `close` double NOT NULL,
    `vol` double NOT NULL,
    PRIMARY KEY(`id`),
    CONSTRAINT `contacts_pk` UNIQUE(`stock_id`, `price_at`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
) ENGINE=INNODB;

CREATE TABLE `user` (
    `id` INT AUTO_INCREMENT,
    `username` CHAR(255) NOT NULL,
    `email` VARCHAR(255) UNIQUE NOT NULL,
    `password` binary(60) NOT NULL,
    PRIMARY KEY(`id`)
) ENGINE=INNODB;

INSERT INTO `user`(`username`, `email`, `password`) VALUES ("test_user", "test@gmail.com", "password");

CREATE TABLE `api` (
    `id` INT AUTO_INCREMENT,
    `type` CHAR(50),
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
    `stock_id` CHAR(50) NOT NULL,
    `price` FLOAT NOT NULL,
    `time` TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
) ENGINE=INNODB;

CREATE TABLE `holiday` (
    `id` INT AUTO_INCREMENT,
    `date` DATE UNIQUE NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB;

-- cross_direction up: TRUE down: FALSE
CREATE TABLE `stock_alert` (
    `id` INT AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `stock_id` CHAR(50) NOT NULL,
    `buy_sell` BOOLEAN NOT NULL,
    `cross_direction` BOOLEAN NOT NULL,
    `price` FLOAT NOT NULL,
    `alert_on` BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
) ENGINE=INNODB;


-- buy_sell buy: TRUE sell: FALSE
CREATE TABLE `trade_log`(
    `id` INT AUTO_INCREMENT,
    `stock_id` CHAR(50) NOT NULL,
    `user_id` INT NOT NULL, 
    `buy_sell` BOOLEAN NOT NULL,
    `time` TIMESTAMP NOT NULL,
    `price` FLOAT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
) ENGINE=INNODB;

CREATE TABLE `cron_hook`(
    `id` INT AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `stock_id` CHAR(50) NOT NULL,
    `price` FLOAT NOT NULL,
    PRIMARY KEY(`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock`(`id`)
) ENGINE=INNODB;