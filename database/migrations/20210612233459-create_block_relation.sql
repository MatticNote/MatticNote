
-- +migrate Up
create table block_relation
(
    block_from uuid not null
        constraint block_from_fk
            references "user"
            on update restrict on delete cascade,
    block_to   uuid not null
        constraint block_to_fk
            references "user"
            on update restrict on delete cascade,
    constraint block_relation_pk
        primary key (block_from, block_to),
    constraint not_same_user_key
        check (block_from <> block_to)
);

-- +migrate Down
drop table if exists block_relation;
