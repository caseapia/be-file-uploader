create table if not exists files_downloads
(
    file_id int not null,
    user_id int not null
);

alter table files_downloads
    add constraint `PRIMARY`
        primary key (file_id);

alter table files_downloads
    add constraint files_downloads_files_id_fk
        foreign key (file_id) references files (id)
            on delete cascade;

alter table files_downloads
    add constraint files_downloads_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade;

