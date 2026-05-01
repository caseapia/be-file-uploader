create table if not exists files_likes
(
    image_id int not null,
    author   int not null
);

create index images_likes_images_id_fk
    on files_likes (image_id);

