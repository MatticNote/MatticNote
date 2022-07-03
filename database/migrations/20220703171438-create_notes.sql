
-- +migrate Up
create table notes
(
    id         char(27)                               not null
        constraint notes_pk
            primary key,
    owner      char(27)
        constraint notes_users_id_fk
            references users
            on update restrict on delete restrict,
    cw         varchar,
    body       text,
    reply_id   char(27)
        constraint notes_reply_id_fk
            references notes
            on update restrict on delete set null,
    retext_id  char(27)
        constraint notes_retext_id_fk
            references notes
            on update restrict on delete set null,
    created_at timestamp with time zone default now() not null
);

create unique index notes_id_uindex
    on notes (id);

-- +migrate Down
drop table notes;
