DROP TABLE users CASCADE;
DROP TABLE post CASCADE;
DROP TABLE comment CASCADE;
DROP TABLE subscription_subscriber;
DROP TABLE user_like_post;
DROP TABLE user_dislike_post;
DROP TABLE user_favorite_post;
DROP TABLE user_like_comment;
DROP TABLE user_dislike_comment;
DROP TABLE schema_migrations;
DROP ROUTINE select_comment(integer);