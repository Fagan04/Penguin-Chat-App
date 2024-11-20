CREATE TABLE users (
                       user_id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(255) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
