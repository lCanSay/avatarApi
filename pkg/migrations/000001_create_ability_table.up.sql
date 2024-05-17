CREATE TABLE IF NOT EXISTS ability 
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    element VARCHAR(255),
    description TEXT,
    image TEXT
);
