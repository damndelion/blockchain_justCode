-- Create the User table
CREATE TABLE Users (
                      id SERIAL PRIMARY KEY,
                      name VARCHAR(255) NOT NULL,
                      email VARCHAR(255) NOT NULL,
                      password VARCHAR(255) NOT NULL,
                      wallet VARCHAR(255) NOT NULL
);

-- Create the UserInfo table
