CREATE TABLE admin (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user"(id) NOT NULL
);