
-- +migrate Up
create table users
(
    id           char(27)                               not null
        constraint users_pk
            primary key,
    username     varchar(64),
    host         varchar
        constraint users_hosts_host_fk
            references hosts
            on update restrict on delete restrict,
    display_name varchar,
    headline     varchar,
    description  text,
    created_at   timestamp with time zone default now() not null,
    is_silence   boolean                  default false not null,
    is_suspend   boolean                  default false not null,
    is_moderator boolean                  default false not null,
    is_admin     boolean                  default false not null,
    deleted_at   timestamp with time zone,
    constraint users_acct_key
        unique (username, host)
);

create unique index users_id_uindex
    on users (id);

create unique index users_acct_uindex
    on users (username, host);

-- +migrate Down
drop table users;
