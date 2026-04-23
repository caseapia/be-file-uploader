create table sessions
(
    id             varchar(255)                         not null
        primary key,
    user_id        int                                  not null,
    ip_address     varchar(45)                          null,
    user_agent     text                                 null,
    is_active      tinyint(1) default 1                 null,
    expires_at     timestamp                            not null,
    created_at     timestamp  default CURRENT_TIMESTAMP null,
    last_active_at timestamp  default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    refresh_hash   varchar(64)                          not null,
    constraint sessions_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade
);

INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('02bdd31e-8f7e-45c7-8d06-832e5554cbe1', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-29 22:18:49', '2026-04-22 22:18:49', '2026-04-22 22:18:49', '71df5955a4afbbfcb93838122a417eb645e519db09980a53f9e9bd3de97c9ee0');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('1d58bd97-3607-4c54-9b42-91d6fe5f2f71', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-24 02:02:41', '2026-04-17 01:59:40', '2026-04-17 01:59:40', '1c916b7d2a9f4a191bd1583115d8abad78bb1b475bacbdf70195acffa8643ac1');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('2b823579-c1fe-4845-b6ea-3d97c6e161be', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-27 18:16:25', '2026-04-20 18:16:25', '2026-04-20 18:16:25', '901539561a121d110f7b2a42690374932eefe8dacbf44a57d9f5cc0aeacb6b72');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('6db34132-799a-46fa-80cd-01def0170a7c', 17, '127.0.0.1', 'Insomnia/be-file-uploader', 1, '2026-04-26 15:48:27', '2026-04-19 15:48:27', '2026-04-19 15:48:27', '92473bb9d59979f1a598526a9bcbf7387dc020552ddbc8d3be0d76d0ba4ef023');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('78770742-0363-4050-8e13-ebc35b36bc3b', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-24 12:05:27', '2026-04-17 12:05:27', '2026-04-17 12:05:27', '72ad485ef2d0e5cbb88ffe3a2f7fbbf0215f891643aec09fd3c19581968658ea');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('b2b49958-2ffb-4f87-bee3-a6bdc3ca53a4', 17, '127.0.0.1', 'Insomnia/be-file-uploader', 1, '2026-04-26 15:48:34', '2026-04-19 15:48:34', '2026-04-19 15:48:34', '1a83bc77ab61467274d46c247cf36630fed0e5b7edd07bf43ae4ab97f6f5cd0c');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('c2952da7-d387-4c56-a2d4-79340768833d', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-29 01:12:50', '2026-04-20 18:26:39', '2026-04-20 18:26:39', '90a652fc8531a00d01c011049ef2337a958889dd5eadbbc5cb4e53d4decf49a8');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('cdbc4f2b-f962-47c3-8c8f-1a9e05fc3e84', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-29 16:24:58', '2026-04-22 16:24:58', '2026-04-22 16:24:58', 'ba6310bd00621a2f76caec84ae7517c73afeb702321d8551fa144c533a3841bd');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('e90b82ea-b5ae-48e0-b717-9dee11261492', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-30 00:42:08', '2026-04-23 00:42:08', '2026-04-23 00:42:08', '0abc239ade9f7aa97459b44d17dfdcf7408b7b9816bfface9978a47389aa9347');
INSERT INTO fileuploader.sessions (id, user_id, ip_address, user_agent, is_active, expires_at, created_at, last_active_at, refresh_hash) VALUES ('efdc5ecd-289c-46bc-bc41-176a47f4c195', 16, '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36 Edg/147.0.0.0', 1, '2026-04-23 21:12:12', '2026-04-16 21:12:12', '2026-04-16 21:12:12', 'c58b98cad1c2cb5de141435cfa7f38e445ff8683c9a75ed2ad90829085cf5b7f');
