DROP TABLE IF EXISTS users;

CREATE TABLE UserInfo (
                          id SERIAL PRIMARY KEY,
                          user_id INT NOT NULL,
                          age INT,
                          phone VARCHAR(20),
                          address VARCHAR(255),
                          country VARCHAR(50),
                          city VARCHAR(50),
                          FOREIGN KEY (user_id) REFERENCES Users(id)  -- Define a foreign key relationship
);

-- Create the UserCredentials table
CREATE TABLE UserCredentials (
                                 id SERIAL PRIMARY KEY,
                                 user_id INT NOT NULL,
                                 card_num VARCHAR(16) NOT NULL,
                                 type VARCHAR(50) NOT NULL,
                                 cvv VARCHAR(3) NOT NULL,
                                 FOREIGN KEY (user_id) REFERENCES Users(id)  -- Define a foreign key relationship
);
