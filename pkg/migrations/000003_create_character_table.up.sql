CREATE TABLE IF NOT EXISTS character (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INTEGER,
    gender VARCHAR(255),
    image VARCHAR(255),
    affiliation_id INTEGER REFERENCES affiliation
);
