CREATE TABLE users
(
    id            bigserial    primary key,
    name          varchar(255) not null,
    surname       varchar(255) not null,
    email         varchar(255) not null,
    password_hash varchar(255) not null
);
