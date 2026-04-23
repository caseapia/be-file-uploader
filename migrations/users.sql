create table users
(
    id           int auto_increment,
    username     varchar(32)                          not null,
    discord_uid  bigint                               null,
    discord_name varchar(32)                          null,
    password     varchar(60)                          null,
    created_at   datetime   default CURRENT_TIMESTAMP not null,
    updated_at   datetime   default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    last_ip      varchar(45)                          not null,
    useragent    text                                 not null,
    invite       int                                  not null,
    register_ip  varchar(45)                          not null,
    upload_limit bigint     default 1073741824        not null,
    cf_ray_id    varchar(32)                          not null,
    used_storage bigint     default 0                 not null,
    is_verified  tinyint(1) default 0                 not null,
    constraint discord_uid
        unique (discord_uid),
    constraint id
        unique (id),
    constraint username
        unique (username)
);