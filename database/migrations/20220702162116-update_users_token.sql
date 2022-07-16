
-- +migrate Up
alter table users_token
    add created_at timestamptz default now() not null;

alter table users_token
    add user_agent varchar;

-- +migrate Down
alter table users_token
    drop column created_at;

alter table users_token
    drop column user_agent;
