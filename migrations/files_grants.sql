create table if not exists files_grants
(
    id         int auto_increment
        primary key,
    file_id    int                                  not null,
    user_id    int                                  not null,
    granted_by int                                  not null,
    granted_at datetime   default CURRENT_TIMESTAMP not null,
    is_owner   tinyint(1) default 0                 not null
);

create index files_grants_file_id_index
    on files_grants (file_id);

create index files_grants_user_id_index
    on files_grants (user_id);

alter table files_grants
    add constraint unique_file_user
        unique (file_id, user_id);

alter table files_grants
    add constraint files_grants_files_id_fk
        foreign key (file_id) references files (id)
            on delete cascade;

alter table files_grants
    add constraint files_grants_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade;

