
-- +migrate Up
create table user_mail_transaction
(
    uuid       uuid                                                            not null
        constraint user_mail_transaction_pk
            primary key
        constraint user_mail_transaction_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    new_email  varchar                                                         not null,
    token      char(32)                                                        not null,
    expired_at timestamp with time zone default (now() + '03:00:00'::interval) not null
);

create unique index user_mail_transaction_new_email_uindex
    on user_mail_transaction (new_email);

create unique index user_mail_transaction_token_uindex
    on user_mail_transaction (token);

create unique index user_mail_transaction_uuid_uindex
    on user_mail_transaction (uuid);

-- +migrate Down
drop table if exists user_mail_transaction;
