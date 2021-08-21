
-- +migrate Up
create table notification
(
    uuid        uuid default gen_random_uuid() not null
        constraint notification_pk
            primary key,
    target_user uuid
        constraint notification_user_uuid_fk_2
            references "user"
            on update restrict on delete cascade,
    from_user   uuid
        constraint notification_user_uuid_fk
            references "user"
            on update restrict on delete set null,
    relate_note uuid
        constraint notification_note_uuid_fk
            references note
            on update restrict on delete cascade,
    type        notification_type not null,
    metadata    jsonb
);

create index notification_target_user_index
    on notification (target_user);

-- +migrate Down
drop table if exists notification;
