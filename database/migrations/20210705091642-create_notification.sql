
-- +migrate Up
create table notification
(
    uuid        uuid              not null
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
    type        integer default 1 not null
        constraint notification_mst_notification_type_notify_id_fk
            references mst_notification_type
            on update cascade on delete set default,
    metadata    jsonb
);

create index notification_target_user_index
    on notification (target_user);

-- +migrate Down
drop table if exists notification;
