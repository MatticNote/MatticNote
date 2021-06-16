
-- +migrate Up
create table mute_relation
(
    mute_from uuid not null
        constraint mute_relation_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    mute_to   uuid not null
        constraint mute_relation_user_uuid_fk_2
            references "user"
            on update restrict on delete cascade,
    constraint mute_relation_pk
        primary key (mute_from, mute_to),
    constraint not_same_user_key
        check (mute_from <> mute_to)
);

-- +migrate Down
drop table if exists mute_relation;
