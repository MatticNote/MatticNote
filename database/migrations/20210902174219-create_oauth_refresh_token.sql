
-- +migrate Up
create table oauth_refresh_token
(
    token      varchar                                                   not null
        constraint oauth_refresh_token_pk
            primary key,
    req_id     varchar,
    client_id  char(32)
        constraint oauth_refresh_token_oauth_client_client_key_fk
            references oauth_client
            on update restrict on delete cascade,
    user_id    uuid
        constraint oauth_refresh_token_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    expires_at timestamp with time zone,
    scopes     character varying[] default ARRAY []::character varying[] not null
);

create unique index oauth_refresh_token_req_id_uindex
    on oauth_refresh_token (req_id);

create unique index oauth_refresh_token_token_uindex
    on oauth_refresh_token (token);

-- +migrate Down
drop table oauth_refresh_token;
