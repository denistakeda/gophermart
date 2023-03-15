create table withdrawals (
    id serial primary key,
    "order" varchar not null,
    sum numeric not null,
    processed_at date not null,
    user_id int not null,

    constraint fk_urser_id
         foreign key(user_id)
         references users(id)
);