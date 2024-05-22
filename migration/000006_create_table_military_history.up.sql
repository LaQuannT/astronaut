CREATE TABLE military_history (
    astronaut_id INT REFERENCES astronaut(id) NOT NULL,
    branch VARCHAR(255) NOT NULL,
    rank VARCHAR(255) NOT NULL,
    retired BOOLEAN DEFAULT FALSE
);
