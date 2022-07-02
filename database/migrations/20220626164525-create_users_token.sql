
-- +migrate Up
create table users_token
(
    id         char(27) not null
        constraint users_token_pk
            primary key,
    token      varchar  not null,
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
