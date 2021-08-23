-- +migrate Up
drop extension if exists pgcrypto;
create extension if not exists pgcrypto;

-- +migrate down
