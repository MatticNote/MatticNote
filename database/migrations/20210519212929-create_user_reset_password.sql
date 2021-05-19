-- +migrate Up
create table user_reset_password
(
    key     char(64)                 not null
        constraint user_reset_password_pk
            primary key,
    target  uuid
        constraint user_email_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    expired timestamp with time zone not null
);

create unique index user_reset_password_key_uindex
    on user_reset_password (key);

create unique index user_reset_password_target_uindex
    on user_reset_password (target);

-- +migrate Down
drop table if exists user_reset_password;
