create table if not exists sessions
(
    id             varchar(36)                            not null,
    user_id        int                                    not null,
    ip_address     varchar(45)                            not null,
    user_agent     text                                   not null,
    is_active      tinyint(1) default 1                   not null,
    created_at     datetime   default CURRENT_TIMESTAMP   not null,
    expires_at     datetime                               not null,
    last_active_at datetime   default CURRENT_TIMESTAMP   not null,
    refresh_hash   varchar(64)                            not null
);

create index sessions_id_index
    on sessions (id);

create index sessions_refresh_hash_index
    on sessions (refresh_hash);

create index sessions_user_id_index
    on sessions (user_id);

alter table sessions
    add constraint `PRIMARY`
        primary key (id);

alter table sessions
    add constraint id
        unique (id);

