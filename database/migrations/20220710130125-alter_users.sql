
-- +migrate Up
alter table users
    add constraint users_admin_moderator_check
        check (NOT (is_moderator IS TRUE AND is_admin IS TRUE));

-- +migrate Down
alter table users
    drop constraint users_admin_moderator_check;
