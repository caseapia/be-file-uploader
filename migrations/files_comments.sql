create table if not exists files_comments
(
    id         int auto_increment
        primary key,
    author     int      null,
    image_id   int      null,
    content    text     null,
    created_at datetime null
);

create index files_comments_files_id_fk
    on files_comments (image_id);

create index files_comments_users_id_fk
    on files_comments (author);

   add constraint files_comments_pk_2
        unique (id);

