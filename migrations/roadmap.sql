create table if not exists roadmap
(
    id         int auto_increment
        primary key,
    title      text                               not null,
    status     tinyint  default 0                 not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    updated_at datetime                           null,
    created_by int                                not null,
    updated_by int                                null
);

create index roadmap_status_index
    on roadmap (status);

alter table roadmap
    add constraint roadmap_pk
        unique (id);

