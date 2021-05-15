
-- +migrate Up
create table user_signature_key
(
    uuid        uuid not null
        constraint signature_key_pk
            primary key
        constraint signature_key_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    public_key  text,
    private_key text
);


-- +migrate Down
drop table if exists user_signature_key;
