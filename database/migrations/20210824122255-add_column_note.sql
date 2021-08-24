
-- +migrate Up
alter table note
    add visibility note_visibility default 'PUBLIC' not null;

-- +migrate Down
alter table note drop column visibility;
