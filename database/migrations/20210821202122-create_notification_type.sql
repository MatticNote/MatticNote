
-- +migrate Up
create type notification_type as enum (
    'REPLY',
    'MENTION',
    'RETEXT',
    'REACTION',
    'FOLLOWED',
    'FOLLOW_REQUEST_RECEIVED',
    'FOLLOW_APPROVED'
);

-- +migrate Down
drop type if exists notification_type;
