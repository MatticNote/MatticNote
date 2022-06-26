
-- +migrate Up
create table users_keypair
(
    id          char(27) not null
        constraint users_private_key_pk
            primary key
        constraint users_private_key_users_id_fk
            references users
            on update restrict on delete restrict,
    private_key bytea,
    public_key  bytea
);

create unique index users_private_key_id_uindex
    on users_keypair (id);

-- +migrate Down
drop table users_keypair;
