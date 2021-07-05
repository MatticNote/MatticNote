
-- +migrate Up
create table mst_notification_type
(
    notify_id int,
    code varchar default 'UNKNOWN'
);

create unique index mst_notification_type_notify_id_uindex
    on mst_notification_type (notify_id);

alter table mst_notification_type
    add constraint mst_notification_type_pk
        primary key (notify_id);

insert into mst_notification_type(notify_id, code) values
    (1, 'UNKNOWN'),
    (2, 'REPLY'),
    (3, 'REACTION'),
    (4, 'RETEXT'),
    (5, 'FOLLOWED'),
    (6, 'FOLLOW_REQUESTED'),
    (7, 'VOTE_CLOSED'),
    (8, 'MENTIONED'),
    (9, 'FOLLOW_APPROVED'),
    (10, 'QUOTE');

-- +migrate Down
drop table if exists mst_notification_type;
