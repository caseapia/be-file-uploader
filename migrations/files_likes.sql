create table if not exists files_likes
(
    image_id int not null,
    author   int not null
);

alter table files_likes
    add constraint `PRIMARY`
        primary key (image_id, author);

alter table files_likes
    add constraint files_likes_files_id_fk
        foreign key (image_id) references files (id)
            on delete cascade;

alter table files_likes
    add constraint files_likes_users_id_fk
        foreign key (author) references users (id)
            on delete cascade;

