CREATE TABLE users
(
    id serial not null primary key,
    username varchar not null unique,
    email varchar not null unique,
    encrypted_password varchar not null,
    "name" varchar,
    surname varchar,
    date_of_birth date,
    date_of_create timestamp
);
CREATE TABLE subscription_subscriber(
    subscription_id int REFERENCES users(id),
    subscriber_id int REFERENCES users(id),
    CONSTRAINT pk_subscribes primary key (subscription_id, subscriber_id)
);


CREATE TABLE post(
    id serial not null primary key,
    text text,
    "object" varchar,
    "views" serial,
    date_of_creation timestamp,
    author_id integer not null,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);
    CREATE TABLE user_favorite_post(
        user_id int REFERENCES users(id),
        post_id int REFERENCES post(id),
        CONSTRAINT pk_favorite_post primary key (user_id, post_id)
    );
    CREATE TABLE user_like_post(
        user_id int REFERENCES users(id),
        post_id int REFERENCES post(id),
        CONSTRAINT pk_like_post primary key (user_id, post_id)
    );
    CREATE TABLE user_dislike_post(
        user_id int REFERENCES users(id),
        post_id int REFERENCES post(id),
        CONSTRAINT pk_dislikes_post primary key (user_id, post_id)
    );


CREATE TABLE comment(
    id serial not null primary key,
    text text,
    date_of_creation timestamp,
    changed boolean default false,
    author_id int not null,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    post_id int not null,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);
    CREATE TABLE parent_child_comment(
        parent_id int REFERENCES comment(id),
        child_id int REFERENCES comment(id),
        CONSTRAINT pk_parent_child_comments primary key (parent_id, child_id)
    );
    CREATE TABLE user_like_comment(
        user_id int REFERENCES users(id),
        comment_id int REFERENCES comment(id),
        CONSTRAINT pk_like_comment PRIMARY KEY (user_id, comment_id)
    );
    CREATE TABLE user_dislike_comment(
        user_id int REFERENCES users(id),
        comment_id int REFERENCES comment(id),
        CONSTRAINT pk_dislike_comment PRIMARY KEY (user_id, comment_id)
    );
