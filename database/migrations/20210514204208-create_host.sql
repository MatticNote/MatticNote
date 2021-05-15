
-- +migrate Up
create table host
(
    host       varchar                                not null
        constraint host_pk
            primary key,
    found_at   timestamp with time zone default now(),
    is_suspend boolean                  default false not null
);

create unique index host_host_uindex
    on host (host);


-- +migrate Down
drop table if exists host;
