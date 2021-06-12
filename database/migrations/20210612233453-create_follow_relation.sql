
-- +migrate Up
create table follow_relation
(
    follow_from uuid not null
        constraint follow_from_fk
            references "user"
            on update restrict on delete cascade,
    follow_to   uuid not null
        constraint follow_to_fk
            references "user"
            on update restrict on delete cascade,
    constraint follow_relation_pk
        primary key (follow_from, follow_to),
    constraint not_same_user_key
        check (follow_from <> follow_to)
);

-- +migrate Down
drop table if exists follow_relation;
