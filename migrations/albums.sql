create table if not exists albums
(
    id         int auto_increment
        primary key,
    name       varchar(64)                          not null,
    created_by int                                  not null,
    created_at datetime   default CURRENT_TIMESTAMP not null,
    updated_at datetime   default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    is_public  tinyint(1) default 0                 not null
)
    charset = utf8mb4;

alter table albums
    add constraint fk_created_by
        foreign key (created_by) references users (id)
            on delete cascade;

