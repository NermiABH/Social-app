CREATE TABLE users(
                      id serial not null primary key,
                      username varchar not null unique,
                      email varchar not null unique,
                      password varchar not null,
                      userpic varchar default 'something url',
                      name varchar,
                      surname varchar,
                      born date,
                      created timestamp default now()
);

CREATE TABLE subscription_subscriber(
                                        subscription_id int REFERENCES users(id) ON DELETE CASCADE,
                                        subscriber_id int REFERENCES users(id) ON DELETE CASCADE,
                                        CONSTRAINT pk_subscribes primary key (subscription_id, subscriber_id)
);

CREATE TABLE post(
                     id serial not null primary key,
                     text text,
                     media varchar[],
                     view_count int,
                     created timestamp default now(),
                     changed boolean default false,
                     author_id integer not null,
                     FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE user_favorite_post(
                                   user_id int REFERENCES users(id) ON DELETE CASCADE,
                                   post_id int REFERENCES post(id) ON DELETE CASCADE,
                                   CONSTRAINT pk_favorite_post primary key (user_id, post_id)
);

CREATE TABLE user_like_post(
                               user_id int REFERENCES users(id) ON DELETE SET NULL,
                               post_id int REFERENCES post(id) ON DELETE CASCADE,
                               CONSTRAINT pk_like_post primary key (user_id, post_id)
);

CREATE TABLE user_dislike_post(
                                  user_id int REFERENCES users(id) ON DELETE SET NULL,
                                  post_id int REFERENCES post(id) ON DELETE CASCADE,
                                  CONSTRAINT pk_dislikes_post primary key (user_id, post_id)
);

CREATE TABLE comment(
                        id serial not null primary key,
                        text text,
                        created timestamp default now(),
                        changed boolean default false,
                        author_id int not null,
                        FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL,
                        post_id int not null,
                        FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
                        parent_id int,
                        FOREIGN KEY (parent_id) REFERENCES comment(id) ON DELETE CASCADE
);
CREATE TABLE user_like_comment(
                                  user_id int REFERENCES users(id) ON DELETE SET NULL,
                                  comment_id int not null REFERENCES comment(id) ON DELETE CASCADE,
                                  CONSTRAINT pk_like_comment PRIMARY KEY (user_id, comment_id)
);

CREATE TABLE user_dislike_comment(
                                     user_id int REFERENCES users(id) ON DELETE SET NULL,
                                     comment_id int REFERENCES comment(id) ON DELETE CASCADE,
                                     CONSTRAINT pk_dislike_comment PRIMARY KEY (user_id, comment_id)
);


-- CREATE OR REPLACE FUNCTION select_comment(int)
--     RETURNS TABLE (id int, author_id int, username varchar, userpic varchar, post_id int, parent_id int,
--         text text, date_of_creation timestamp, changed bool, likes bigint, dislikes bigint)
--     LANGUAGE sql
-- AS $function$
-- WITH RECURSIVE r AS (
--     (SELECT id
--      FROM comment
--      WHERE id = $1)
--     UNION ALL
--     SELECT c.id
--     FROM r
--         JOIN comment c ON c.parent_id = r.id
--     WHERE r.id < (SELECT id
--         FROM comment
--         ORDER BY date_of_creation DESC
--         LIMIT 1)
-- )
-- SELECT  c.id, c.author_id, u.username, u.userpic, c.post_id, c.parent_id,
--         c.text, c.date_of_creation, c.changed,
--         COUNT(ulc) as likes, COUNT(udc) as dislikes
-- FROM r rr
--          LEFT JOIN comment c ON rr.id = c.id
--          LEFT JOIN users u on c.author_id = u.id
--          LEFT JOIN user_like_comment ulc on rr.id = ulc.comment_id
--          LEFT JOIN user_dislike_comment udc on rr.id = udc.comment_id
-- GROUP BY c.author_id, u.username, u.userpic, c.post_id, c.parent_id, c.text, c.date_of_creation, c.changed, c.id
-- ORDER BY c.date_of_creation
-- $function$;
--

INSERT INTO users (username, email, password) VALUES ('test1', 'test1@gmail.com', 'pusinu48') returning id;
INSERT INTO post (author_id, text) VALUES (5, 'fsadfasdf');
SELECT id FROM post OFFSET 5 LIMIT 5;