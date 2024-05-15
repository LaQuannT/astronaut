CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE api_key (
    key UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id INT REFERENCES "user"(id) NOT NULL
);