create table customer (
  id serial primary key,
  first_name varchar not null,
  last_name varchar not null,
  birth_date date,
  creation_time timestamptz not null default now()
);

create table widget (
  id bigserial primary key,
  name varchar not null,
  weight smallint not null
);
