
-- +migrate Up
create table user_auth
(
    id       char(27)    not null
        constraint user_password_pk
            primary key
        constraint user_password_user_id_fk
            references "user"
            on update restrict on delete restrict,
    password varchar(64) not null
);

create unique index user_auth_id_uindex
    on user_auth (id);

-- +migrate Down
drop table user_auth;
