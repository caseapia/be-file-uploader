create table if not exists albums
(
    id         int auto_increment
        primary key,
    name       varchar(64)                          not null,
    created_by int                                  not null,
    created_at datetime   default CURRENT_TIMESTAMP not null,
    updated_at datetime                             not null,
    is_public  tinyint(1) default 0                 not null
);

alter table albums
    add constraint albums_pk_2
        unique (id);

