CREATE TABLE astronaut_mission (
    astronaut_id INT REFERENCES astronaut(id) NOT NULL,
    mission_id INT REFERENCES mission(id) NOT NULL
);