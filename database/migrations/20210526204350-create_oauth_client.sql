
-- +migrate Up
create table oauth_client
(
    client_key    char(32)                                                  not null
        constraint oauth_client_pk
            primary key,
    client_secret bytea                                                     not null,
    name          varchar(64)                                               not null,
    client_owner  uuid
        constraint oauth_client_user_uuid_fk
            references "user"
            on update restrict on delete set null,
    redirect_uris character varying[] default ARRAY []::character varying[] not null,
    scopes        character varying[] default ARRAY []::character varying[] not null
);

create unique index oauth_client_client_key_uindex
    on oauth_client (client_key);

-- +migrate Down
drop table if exists oauth_client;
