create table if not exists user_roles
(
    user_id int not null,
    role_id int not null
)
    comment 'роли пользователей';

alter table user_roles
    add constraint user_roles_users_id_fk
        foreign key (user_id) references users (id)
            on update cascade on delete cascade;

