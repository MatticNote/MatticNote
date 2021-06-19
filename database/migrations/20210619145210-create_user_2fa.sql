
-- +migrate Up
create table user_2fa
(
    uuid        uuid                  not null
        constraint user_2fa_pk
            primary key
        constraint user_2fa_user_uuid_fk
            references "user"
            on update cascade on delete cascade,
    is_enable   boolean default false not null,
    secret_code char(32)              not null,
    backup_code jsonb   default '[]'::jsonb
);

-- +migrate Down
drop table if exists user_2fa;
