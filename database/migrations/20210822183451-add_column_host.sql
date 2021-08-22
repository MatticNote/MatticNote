
-- +migrate Up
alter table host
    add shared_inbox varchar;

-- +migrate Down
alter table host drop column shared_inbox;
