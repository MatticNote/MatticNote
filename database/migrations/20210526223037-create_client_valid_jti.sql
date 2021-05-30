
-- +migrate Up
create table oauth_valid_jti
(
    jti        varchar                                                         not null
        constraint oauth_valid_jti_pk
            primary key,
    expired_at timestamp with time zone default (now() + '01:00:00'::interval) not null
);

create unique index oauth_valid_jti_jti_uindex
    on oauth_valid_jti (jti);

-- +migrate Down
drop table if exists oauth_valid_jti;
