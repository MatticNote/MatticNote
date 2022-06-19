
-- +migrate Up
create table note
(
    id         char(27) not null
        constraint note_pk
            primary key,
    owner      char(27)
        constraint note_user_id_fk
            references "user"
            on update restrict on delete restrict,
    cw         varchar,
    body       text,
    created_at timestamp with time zone default now()
);

create unique index note_id_uindex
    on note (id);

-- +migrate Down
drop table note;
