CREATE TABLE users (
    id varchar(50) PRIMARY KEY,
    name varchar(150) NOT NULL CHECK (char_length(name) > 0),
    password bytea NOT NULL,
    picture varchar(24)
);

CREATE TABLE usr_del_confirm_codes (
    userId varchar(50) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    code varchar(36) NOT NULL
);

CREATE TABLE deleted_users (
    id varchar(50) PRIMARY KEY
);