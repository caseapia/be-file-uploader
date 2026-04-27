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

create table if not exists files
(
    id            int auto_increment
        primary key,
    r2_key        varchar(512)         not null,
    url           text                 not null,
    original_name varchar(255)         null,
    mime_type     varchar(64)          null,
    size          bigint               null,
    uploaded_by   int                  not null,
    is_private    tinyint(1) default 0 not null,
    album_id      int                  null,
    downloads     int                  null
);

create index idx_images_uploaded_by
    on files (uploaded_by);

alter table files
    add constraint r2_key
        unique (r2_key);

alter table files
    add constraint images_albums_id_fk
        foreign key (album_id) references albums (id)
            on delete cascade;

create table if not exists files_comments
(
    id         int auto_increment
        primary key,
    author     int                                not null,
    image_id   int                                not null,
    content    text                               not null,
    created_at datetime default CURRENT_TIMESTAMP not null
);

alter table files_comments
    add constraint files_comments_pk_2
        unique (id);

create table if not exists files_downloads
(
    file_id int not null,
    user_id int not null
);

alter table files_downloads
    add constraint `PRIMARY`
        primary key (file_id);

create table if not exists files_likes
(
    image_id int not null,
    author   int not null
);

create index images_likes_images_id_fk
    on files_likes (image_id);

create table if not exists invites
(
    id         int auto_increment
        primary key,
    code       varchar(6)                           not null,
    created_by int                                  not null,
    is_active  tinyint(1) default 1                 null,
    created_at timestamp  default CURRENT_TIMESTAMP null,
    used_by    int                                  null
);

create index fk_inviter
    on invites (created_by);

alter table invites
    add constraint code
        unique (code);

create table if not exists roles
(
    id          int auto_increment
        primary key,
    name        varchar(100)                         not null,
    permissions json                                 not null,
    is_system   tinyint(1) default 0                 null,
    created_at  timestamp  default CURRENT_TIMESTAMP null,
    created_by  int                                  not null,
    color       varchar(9) default '#ffffff'         not null
);

alter table roles
    add constraint id
        unique (id);

alter table roles
    add constraint name
        unique (name);

alter table roles
    add constraint roles_pk
        unique (id);

create table if not exists user_roles
(
    user_id int not null,
    role_id int not null
)
    comment 'роли пользователей';

create table if not exists users
(
    id           int,
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
    locale       varchar(2)                           not null
);

alter table users
    add constraint discord_uid
        unique (discord_uid);

alter table users
    add constraint id
        unique (id);

alter table files
    add constraint fk_images_user
        foreign key (uploaded_by) references users (id)
            on delete cascade;

alter table user_roles
    add constraint user_roles_users_id_fk
        foreign key (user_id) references users (id)
            on update cascade on delete cascade;

alter table users
    add constraint username
        unique (username);

alter table users
    modify id int auto_increment;


