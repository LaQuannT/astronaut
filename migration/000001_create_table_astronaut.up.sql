CREATE TABLE astronaut (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL ,
    last_name VARCHAR(255) NOT NULL,
    gender CHAR(1) CHECK ( gender = 'F' OR gender = 'M' ),
    birth_date DATE NOT NULL,
    birth_place VARCHAR(255) NOT NULL
);