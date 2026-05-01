CREATE DATABASE IF NOT EXISTS sprinter_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE sprinter_db;

CREATE TABLE users
(
    id                INT PRIMARY KEY AUTO_INCREMENT,
    name              VARCHAR(50)    NOT NULL,
    email             VARCHAR(50)    NOT NULL,
    password          VARCHAR(255)   NOT NULL,
    username          VARCHAR(50)    NOT NULL DEFAULT '',
    biography         VARCHAR(255)   NOT NULL DEFAULT '',
    image_url         VARCHAR(255)   NOT NULL DEFAULT '',
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
CREATE TABLE activities (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            uuid CHAR(36) NOT NULL UNIQUE,
                            user_id INT NOT NULL,
                            type INT NOT NULL, -- 1=Cycling, 2=Running, 3=Walking

                            start_time DATETIME(6),
                            end_time DATETIME(6),

                            created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
                            modified_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

                            INDEX idx_user_id (user_id)
);
CREATE TABLE points (
                        id INT AUTO_INCREMENT PRIMARY KEY,
                        uuid CHAR(36) NOT NULL UNIQUE,
                        activity_id INT NOT NULL,

                        latitude DOUBLE NOT NULL,
                        longitude DOUBLE NOT NULL,

                        created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
                        modified_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

                        CONSTRAINT fk_points_activity
                            FOREIGN KEY (activity_id)
                                REFERENCES activities(id)
                                ON DELETE CASCADE,

                        INDEX idx_activity_id (activity_id)
);