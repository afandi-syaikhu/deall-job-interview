insert into users
values
    ( 1, 'admin@email.com', MD5('admin'), 'admin', now(), now()),
    ( 2, 'user@email.com', MD5('user'), 'user', now(), now());

ALTER SEQUENCE users_id_seq RESTART WITH 3;