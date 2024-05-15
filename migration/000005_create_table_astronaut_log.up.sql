CREATE TYPE status AS ENUM ('active', 'retired', 'management', 'deceased');

CREATE TABLE astronaut_log (
    astronaut_id INT REFERENCES astronaut(id) NOT NULL,
    space_flights INT DEFAULT 0,
    space_flight_hrs INT DEFAULT 0,
    space_walks INT DEFAULT 0,
    space_walk_hrs INT DEFAULT 0,
    status status NOT NULL,
    death_mission INT REFERENCES mission(id),
    death_date DATE
);