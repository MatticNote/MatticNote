
-- +migrate Up
create table user_session
(
    token       char(128) not null
        constraint user_session_pk
            primary key,
    user_id     char(27)
        constraint user_session_user_id_fk
            references "user"
            on update restrict on delete restrict,
    expired_at  timestamp with time zone,
    issued_from inet
);

create unique index user_session_token_uindex
    on user_session (token);

-- +migrate Down
drop table user_session;
