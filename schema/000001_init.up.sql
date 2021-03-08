CREATE TABLE users
(
    id            bigserial                primary key,
    email         varchar(255)             not null,
    username      varchar(255)             not null UNIQUE,
    first_name    varchar(255)             not null,
    last_name     varchar(255)             not null,
    password      varchar(255)             not null,
    created_at    timestamp with time zone not null default timezone('UTC'::text, now()),
    updated_at    timestamp with time zone not null default timezone('UTC'::text, now())
);
