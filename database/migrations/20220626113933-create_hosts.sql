
-- +migrate Up
create table hosts
(
    host        varchar                                not null
        constraint hosts_pk
            primary key,
    name        varchar,
    description text,
    is_suspend  boolean                  default false not null,
    found_at    timestamp with time zone default now() not null
);

-- +migrate Down
drop table hosts;
