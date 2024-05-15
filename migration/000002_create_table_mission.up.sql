CREATE TABLE mission (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    alias VARCHAR(255),
    date_of_mission DATE NOT NULL,
    successful BOOLEAN DEFAULT TRUE NOT NULL
);