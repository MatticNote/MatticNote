
-- +migrate Up
alter table follow_relation
    add is_pending bool default false not null;

-- +migrate Down
alter table follow_relation drop column is_pending;
