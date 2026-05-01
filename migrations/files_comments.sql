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

