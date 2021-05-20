
-- +migrate Up
alter table user_reset_password alter column target set not null;

-- +migrate Down
alter table user_reset_password alter column target drop not null;
