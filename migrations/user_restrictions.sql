create table if not exists user_restrictions
(
    id           int auto_increment
        primary key,
    user_id      int                                                                           not null,
    moderator_id int                                                                           not null,
    created_at   datetime                                            default CURRENT_TIMESTAMP not null,
    unban_at     datetime                                                                      null,
    reason       varchar(255)                                                                  not null,
    status       enum ('banned', 'unbanned_by_admin', 'ban_expired') default 'banned'          not null,
    unbanned_by  int                                                                           null,
    type         enum ('account', 'upload', 'like', 'comment')       default 'account'         not null
)
    charset = utf8mb4;

create index bans_type_index
    on user_restrictions (status);

create index user_restrictions_type_index
    on user_restrictions (type);

alter table user_restrictions
    add constraint bans_users_id1_fk
        foreign key (moderator_id) references users (id)
            on delete cascade;

alter table user_restrictions
    add constraint bans_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade;

