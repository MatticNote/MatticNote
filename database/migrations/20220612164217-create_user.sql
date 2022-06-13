
-- +migrate Up
create table "user"
(
    id           char(27)                               not null
        constraint user_pk
            primary key,
    username     varchar(32),
    host         varchar(255)
        constraint user_fediverse_host_fk
            references fediverse
            on update restrict on delete restrict,
    display_name varchar,
    headline     varchar,
    description  text,
    created_at   timestamp with time zone default now(),
    is_silence   boolean                  default false not null,
    is_suspend   boolean                  default false not null,
    is_active    boolean                  default true  not null,
    is_moderator boolean                  default false not null,
    is_admin     boolean                  default false not null,
    constraint acct_key
        unique (username, host)
);

create unique index user_id_uindex
    on "user" (id);

-- +migrate Down
drop table "user";
