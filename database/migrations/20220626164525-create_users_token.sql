
-- +migrate Up
create table users_token
(
    token      varchar not null
        constraint users_token_pk
            primary key,
    user_id    char(27)
        constraint users_token_users_id_fk
            references users
            on update restrict on delete restrict,
    expired_at timestamp with time zone,
    ip         inet
);

create unique index users_token_token_uindex
    on users_token (token);

-- +migrate Down
drop table users_token;
