
-- +migrate Up
create table users_auth
(
    id       char(27) not null
        constraint users_auth_pk
            primary key
        constraint users_auth_users_id_fk
            references users
            on update restrict on delete restrict,
    password bytea
);

create unique index users_auth_id_uindex
    on users_auth (id);

-- +migrate Down
drop table users_auth;
