
-- +migrate Up
create table user_email
(
    id       char(27)              not null
        constraint user_email_pk
            primary key
        constraint user_email_user_id_fk
            references "user"
            on update restrict on delete restrict,
    email    varchar(255),
    verified boolean default false not null
);

create unique index user_email_email_uindex
    on user_email (email);

-- +migrate Down
drop table user_email;
