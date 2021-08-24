
-- +migrate Up
create type note_visibility as enum (
    'PUBLIC',
    'UNLISTED',
    'FOLLOWER',
    'DIRECT'
);

-- +migrate Down
drop type if exists note_visibility;
