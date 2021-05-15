
-- +migrate Up
create table "user"
(
    uuid            uuid                     default gen_random_uuid() not null
        constraint user_pk
            primary key,
    username        varchar(32)                                        not null,
    host            varchar
        constraint user_host_host_fk
            references host
            on update restrict on delete restrict,
    display_name    varchar,
    summary         text,
    created_at      timestamp with time zone default now(),
    updated_at      timestamp with time zone default now(),
    is_active       boolean                  default true              not null,
    is_silence      boolean                  default false             not null,
    is_suspend      boolean                  default false             not null,
    accept_manually boolean                  default false             not null,
    is_superuser    boolean                  default false             not null,
    avatar_uuid     uuid, -- Foreign key will be set in the future
    header_uuid     uuid, -- Foreign key will be set in the future
    is_bot          boolean                  default false             not null,
    constraint acct_pk
        unique (username, host)
);

-- +migrate Down
drop table if exists "user";
