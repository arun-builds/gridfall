-- +goose Up
CREATE TYPE account_type AS ENUM (
    'guest',
    'registered'
);

CREATE TYPE user_role AS ENUM (
    'player',
    'admin'
);

create table users(
    id uuid primary key default gen_random_uuid(),
    name varchar(50) not null,
    email varchar(255) unique,
    password_hash text,
    type account_type not null,
    role user_role NOT NULL DEFAULT 'player',
    created_at timestamptz NOT NULL DEFAULT now()
);



-- +goose Down
DROP TABLE users;
DROP TYPE user_role;
DROP TYPE account_type;
