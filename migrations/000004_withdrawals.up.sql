create table withdrawals (
    id serial primary key,
    order_number varchar not null,
    sum bigint not null,
    processed_at date not null,
    user_id int not null,

    constraint fk_urser_id
         foreign key(user_id)
         references users(id)
);