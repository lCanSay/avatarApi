CREATE TABLE IF NOT EXISTS affiliation (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image TEXT
);
