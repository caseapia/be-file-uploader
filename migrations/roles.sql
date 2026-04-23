create table roles
(
    id          int auto_increment
        primary key,
    name        varchar(100)                         not null,
    permissions json                                 not null,
    is_system   tinyint(1) default 0                 null,
    created_at  timestamp  default CURRENT_TIMESTAMP null,
    created_by  int                                  not null,
    color       varchar(9) default '#ffffff'         not null,
    constraint id
        unique (id),
    constraint name
        unique (name),
    constraint roles_pk
        unique (id)
);