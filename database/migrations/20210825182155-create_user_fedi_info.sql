
-- +migrate Up
create table user_fedi_info
(
    uuid     uuid not null
        constraint user_fedi_info_pk
            primary key
        constraint user_fedi_info_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    inbox    varchar,
    outbox   varchar,
    featured varchar
);

create unique index user_fedi_info_uuid_uindex
    on user_fedi_info (uuid);


-- +migrate Down
drop table if exists user_fedi_info;
