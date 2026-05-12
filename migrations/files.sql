create table if not exists files
(
    id            int auto_increment
        primary key,
    r2_key        varchar(512)                       not null,
    url           text                               not null,
    original_name varchar(255)                       not null,
    mime_type     varchar(64)                        not null,
    size          bigint                             not null,
    uploaded_by   int                                not null,
    is_private    tinyint  default 0                 not null,
    album_id      int                                null,
    downloads     int      default 0                 not null,
    created_at    datetime default CURRENT_TIMESTAMP not null
);

create index idx_images_uploaded_by
    on files (uploaded_by);

create index images_albums_id_fk_2
    on files (album_id);

alter table files
    add constraint r2_key
        unique (r2_key);

alter table files
    add constraint fk_images_user
        foreign key (uploaded_by) references users (id)
            on delete cascade;

alter table files
    add constraint images_albums_id_fk
        foreign key (album_id) references albums (id)
            on delete cascade;

