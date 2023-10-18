CREATE TABLE user_info (
                          id SERIAL PRIMARY KEY,
                          user_id INT NOT NULL,
                          age INT,
                          phone VARCHAR(20),
                          address VARCHAR(255),
                          country VARCHAR(50),
                          city VARCHAR(50),
                          FOREIGN KEY (user_id) REFERENCES Users(id)  -- Define a foreign key relationship
);
