
-- +migrate Up
alter table notification alter column uuid set default gen_random_uuid();

-- +migrate Down
alter table notification alter column uuid drop default;
