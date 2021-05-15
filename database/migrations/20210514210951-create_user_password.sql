
-- +migrate Up
create table user_password
(
    uuid     uuid not null
        constraint user_password_pk
            primary key
        constraint user_password_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    password bytea
);

create unique index user_password_uuid_uindex
    on user_password (uuid);


-- +migrate Down
drop table if exists user_password;
