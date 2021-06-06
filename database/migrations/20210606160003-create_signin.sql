
-- +migrate Up
create table signin
(
    uuid        uuid                     default gen_random_uuid() not null
        constraint signin_pk
            primary key,
    tried_at    timestamp with time zone default now()             not null,
    target_user uuid
        constraint signin_user_uuid_fk
            references "user"
            on update restrict on delete set null,
    is_success  boolean                  default false             not null,
    from_ip     inet
);

create index signin_target_user_index
    on signin (target_user);

-- +migrate Down
drop table if exists signin;
