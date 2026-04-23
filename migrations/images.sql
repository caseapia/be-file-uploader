create table images
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
    constraint r2_key
        unique (r2_key),
    constraint fk_images_user
        foreign key (uploaded_by) references users (id)
            on delete cascade,
    constraint images_albums_id_fk
        foreign key (album_id) references albums (id)
            on delete cascade
);

create index idx_images_uploaded_by
    on images (uploaded_by);

