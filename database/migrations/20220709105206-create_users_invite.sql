
-- +migrate Up
create table users_invite
(
    id         char(27) not null
        constraint users_invite_pk
            primary key,
    owner      char(27)
        constraint users_invite_users_id_fk
            references users
            on update restrict on delete restrict,
    code       varchar  not null,
    count      integer,
    expired_at timestamp with time zone
);

create unique index users_invite_code_uindex
    on users_invite (code);

create unique index users_invite_id_uindex
    on users_invite (id);

-- +migrate Down
drop table users_invite;
