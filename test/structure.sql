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
