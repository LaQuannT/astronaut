CREATE TABLE astronaut_undergrad_major (
    astronaut_id INT REFERENCES astronaut(id) NOT NULL,
    major_id INT REFERENCES major(id) NOT NULL
);