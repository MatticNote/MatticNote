
-- +migrate Up
alter table "user"
    add ap_id varchar;

create unique index user_ap_id_uindex
    on "user" (ap_id);

-- +migrate Down
drop index user_ap_id_uindex;
alter table "user" drop column ap_id;
