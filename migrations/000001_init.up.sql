create table users (
    id serial primary key,
    login varchar(255) unique not null,
    password varchar(255) not null,
    created_at timestamp not null
)