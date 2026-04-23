create table images_likes
(
    image_id int not null,
    author   int not null,
    constraint images_likes_images_id_fk
        foreign key (image_id) references images (id)
            on delete cascade
);

