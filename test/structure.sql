drop table if exists customer;
create table customer (
  id serial primary key,
  first_name varchar not null,
  last_name varchar not null,
  birth_date date,
  creation_time timestamptz not null default now()
);

drop table if exists widget;
create table widget (
  id bigserial primary key,
  name varchar not null,
  weight smallint not null
);

drop table if exists part;
create table part (
  code varchar primary key,
  description text not null
);

drop table if exists semester;
create table semester (
  year smallint,
  season varchar,
  description text not null,
  primary key(year, season)
);

drop table if exists blob;
create table blob (
  id serial primary key,
  payload bytea not null
);
