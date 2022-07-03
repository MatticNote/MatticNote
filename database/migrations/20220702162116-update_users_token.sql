
-- +migrate Up
alter table users_token
    add created_at timestamptz default now() not null;

-- +migrate Down
alter table users_token
    drop column created_at;
