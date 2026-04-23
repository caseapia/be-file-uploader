create table user_roles
(
    user_id int not null,
    role_id int not null,
    constraint user_roles_users_id_fk
        foreign key (user_id) references users (id)
            on update cascade on delete cascade
)
    comment 'роли пользователей';
