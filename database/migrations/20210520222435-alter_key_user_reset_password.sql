
-- +migrate Up
alter table user_reset_password
    add constraint user_reset_password_target
        unique (target);

-- +migrate Down
alter table user_reset_password drop constraint user_reset_password_target;
