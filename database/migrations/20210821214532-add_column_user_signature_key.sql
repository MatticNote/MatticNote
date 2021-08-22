
-- +migrate Up
alter table user_signature_key
    add key_id varchar;

create unique index user_signature_key_key_id_uindex
    on user_signature_key (key_id);

-- +migrate Down
drop index user_signature_key_key_id_uindex;

alter table user_signature_key drop column key_id;
