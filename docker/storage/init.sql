CREATE TABLE users (
    id varchar(50) PRIMARY KEY,
    name varchar(150) NOT NULL CHECK (char_length(name) > 0),
    password bytea NOT NULL,
    picture varchar
);