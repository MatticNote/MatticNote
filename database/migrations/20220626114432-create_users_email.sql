
-- +migrate Up
create table users_email
(
    id          char(27)              not null
        constraint users_email_pk
            primary key
        constraint users_email_users_id_fk
            references users
            on update restrict on delete restrict,
    email       varchar,
    is_verified boolean default false not null
);

create unique index users_email_id_uindex
    on users_email (id);

create index users_email_email_index
    on users_email (email);

-- +migrate Down
drop table users_email;
