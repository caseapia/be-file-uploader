create table roadmap
(
    id          int auto_increment
        primary key,
    title       int                                not null,
    tasks       json                               not null,
    created_at  datetime default CURRENT_TIMESTAMP not null,
    created_by  int                                not null,
    finished_at datetime                           null,
    constraint id
        unique (id),
    constraint roadmap_pk_2
        unique (id)
);

