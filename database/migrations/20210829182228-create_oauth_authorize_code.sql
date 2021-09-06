
-- +migrate Up
create table oauth_authorize_code
(
    code       varchar                                                         not null
        constraint oauth_authorize_code_pk
            primary key,
    expires_at timestamp with time zone default (now() + '00:15:00'::interval) not null,
    scopes     character varying[]      default ARRAY []::character varying[]  not null,
    client_id  char(32)
        constraint oauth_authorize_code_oauth_client_client_key_fk
            references oauth_client
            on update restrict on delete cascade,
    user_id    uuid
        constraint oauth_authorize_code_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    is_active  boolean                  default true                           not null
);

create unique index oauth_authorize_code_code_uindex
    on oauth_authorize_code (code);

-- +migrate Down
drop table oauth_authorize_code;
