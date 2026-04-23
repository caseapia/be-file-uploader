create table invites
(
    id         int auto_increment
        primary key,
    code       varchar(6)                           not null,
    created_by int                                  not null,
    is_active  tinyint(1) default 1                 null,
    created_at timestamp  default CURRENT_TIMESTAMP null,
    used_by    int                                  null,
    constraint code
        unique (code)
);

create index fk_inviter
    on invites (created_by);
