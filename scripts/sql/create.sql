CREATE DATABASE IF NOT EXISTS sprinter_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE TABLE users
(
    id                INT PRIMARY KEY AUTO_INCREMENT,
    name              VARCHAR(50)    NOT NULL,
    username          VARCHAR(50)    NOT NULL,
    biography         VARCHAR(255)   NOT NULL,
    image_url         VARCHAR(255)   NOT NULL,
    email             VARCHAR(50)    NOT NULL,
    password          VARCHAR(255)   NOT NULL,
    carbo_coins       INT            NOT NULL DEFAULT 0,
    carbon            DECIMAL(10, 2) NOT NULL DEFAULT 0,
    traveled_distance DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status_code       INT            NOT NULL DEFAULT 0,
    created_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at       TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_users_email (email),
    INDEX idx_user_email (email),
    INDEX idx_user_carbo_coins (carbo_coins),
    INDEX idx_user_carbon (carbon),
    INDEX idx_user_status_code (status_code)
);

CREATE TABLE products
(
    id          INT PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(50)  NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    price       INT          NOT NULL,
    stock       INT          NOT NULL DEFAULT 0,
    image_url   VARCHAR(255) NOT NULL DEFAULT '',
    status_code INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_product_status_code (status_code)
);