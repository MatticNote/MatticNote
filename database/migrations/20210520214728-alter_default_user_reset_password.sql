
-- +migrate Up
alter table user_reset_password alter column expired set default now() + interval '30 minutes';

-- +migrate Down
alter table user_reset_password alter column expired drop default;
