CREATE TABLE astronaut_alma_mater (
    astronaut_id INT REFERENCES astronaut(id) NOT NULL,
    alma_mater_id INT REFERENCES alma_mater(id) NOT NULL
);