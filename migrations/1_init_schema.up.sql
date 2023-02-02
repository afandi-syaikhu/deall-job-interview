CREATE TYPE user_role AS ENUM('admin', 'user');

CREATE TABLE users (
    id bigserial PRIMARY KEY,
    username varchar(50) NOT NULL,
    password varchar(255) NOT NULL,
    role user_role NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users
    ADD CONSTRAINT unique_username UNIQUE (username);