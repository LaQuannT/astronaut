CREATE TABLE military_history (
    astronaut_id INT REFERENCES astronaut(id) NOT NULL,
    branch VARCHAR(255),
    rank VARCHAR(255),
    retired BOOLEAN DEFAULT FALSE
);