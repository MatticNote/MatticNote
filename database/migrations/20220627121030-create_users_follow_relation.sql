
-- +migrate Up
create table users_follow_relation
(
    id          char(27)              not null
        constraint users_follow_relation_pk
            primary key,
    from_follow char(27)              not null
        constraint users_follow_relation_users_id_fk
            references users
            on update restrict on delete restrict,
    to_follow   char(27)              not null
        constraint users_follow_relation_users_id_fk_2
            references users
            on update restrict on delete restrict,
    is_active   boolean default false not null,
    constraint users_follow_relation_key
        unique (from_follow, to_follow),
    constraint users_follow_relation_do_not_self_follow
        check (from_follow <> to_follow)
);

create unique index users_follow_relation_key_uindex
    on users_follow_relation (from_follow, to_follow);

-- +migrate Down
drop table users_follow_relation;
