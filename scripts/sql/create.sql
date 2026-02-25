CREATE DATABASE IF NOT EXISTS sprinter_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE sprinter_db;

CREATE TABLE users
(
    id                INT PRIMARY KEY AUTO_INCREMENT PRIMARY KEY,
    name              VARCHAR(50)    NOT NULL,
    email             VARCHAR(50)    NOT NULL,
    password          VARCHAR(255)   NOT NULL,
    carbo_coins       INT            NOT NULL DEFAULT 0,
    carbon            DECIMAL(10, 2) NOT NULL DEFAULT 0,
    traveled_distance DECIMAL(10, 2) NOT NULL DEFAULT 0,
    UNIQUE KEY uq_users_email (email)
);
