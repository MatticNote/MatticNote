
-- +migrate Up
create table fediverse
(
    host        varchar(255)                           not null
        constraint fediverse_pk
            primary key,
    name        varchar(255),
    description text,
    blocked     boolean                  default false not null,
    found_at    timestamp with time zone default now() not null
);

create unique index fediverse_host_uindex
    on fediverse (host);

-- +migrate Down
drop table fediverse;
