CREATE TABLE chat_members (
                              chat_member_id INT AUTO_INCREMENT PRIMARY KEY,
                              chat_id INT NOT NULL,
                              user_id INT NOT NULL,
                              joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              FOREIGN KEY (chat_id) REFERENCES chats(chat_id) ON DELETE CASCADE,
                              FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
