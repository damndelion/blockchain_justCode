CREATE TABLE user_credentials (
                                 id SERIAL PRIMARY KEY,
                                 user_id INT NOT NULL,
                                 card_num VARCHAR(16) NOT NULL,
                                 type VARCHAR(50) NOT NULL,
                                 cvv VARCHAR(3) NOT NULL,
                                 FOREIGN KEY (user_id) REFERENCES Users(id)  -- Define a foreign key relationship
);
