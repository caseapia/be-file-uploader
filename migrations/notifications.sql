create table if not exists notifications
(
    id         int auto_increment
        primary key,
    user_id    int                                  not null,
    content    text                                 not null,
    created_at datetime   default CURRENT_TIMESTAMP not null,
    is_readed  tinyint(1) default 0                 not null
);

create index notifications_user_id_index
    on notifications (user_id desc);

alter table notifications
    add constraint notifications_pk_2
        unique (id);

