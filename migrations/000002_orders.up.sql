create table orders (
    id serial primary key,
    user_id int not null,
    order_number int not null,
    status varchar not null,
    created_at timestamp not null,
    updated_at timestamp not null,

    constraint fk_user_id
        foreign key(user_id)
        references users(id),

    constraint unique_user_id_and_order_number
        unique (user_id, order_number)
);