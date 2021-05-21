
-- +migrate Up
create table note
(
    uuid        uuid                     default gen_random_uuid() not null
        constraint note_pk
            primary key,
    author      uuid                                               not null
        constraint note_user_uuid_fk
            references "user"
            on update restrict on delete cascade,
    created_at  timestamp with time zone default now(),
    cw          text,
    body        text,
    reply_uuid  uuid
        constraint note_reply_uuid_fkey
            references note
            on update restrict on delete set null,
    retext_uuid uuid
        constraint note_retext_uuid_fkey
            references note
            on update restrict on delete set null,
    local_only  boolean                  default false             not null
);

create index note_author_index
    on note (author);

-- +migrate Down
drop table if exists note;
