
-- +migrate Up
create table user_mail
(
    uuid        uuid                  not null
        constraint user_mail_verified_pk
            primary key
        constraint user_mail_verified_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    email       varchar               not null,
    is_verified boolean default false not null
);

create unique index user_mail_verified_uuid_uindex
    on user_mail (uuid);

create unique index user_mail_email_uindex
    on user_mail (email);


-- +migrate Down
drop table if exists user_mail;
