
-- +migrate Up
create table oauth_access_token
(
    token     varchar not null
        constraint oauth_access_token_pk
            primary key,
    req_id    varchar,
    expires   timestamp with time zone default (now() + '01:00:00'::interval),
    scopes    character varying[]      default ARRAY []::character varying[],
    client_id char(32)
        constraint oauth_access_token_oauth_client_client_key_fk
            references oauth_client
            on update restrict on delete cascade,
    user_id   uuid
        constraint oauth_access_token_user_uuid_fk
            references "user"
);

create unique index oauth_access_token_req_id_uindex
    on oauth_access_token (req_id);

create unique index oauth_access_token_token_uindex
    on oauth_access_token (token);


-- +migrate Down
drop table oauth_access_token;
