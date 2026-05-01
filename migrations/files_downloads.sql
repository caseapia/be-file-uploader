create table if not exists files_downloads
(
    file_id int not null,
    user_id int not null
);

alter table files_downloads
    add constraint `PRIMARY`
        primary key (file_id);

