CREATE TABLE collections (
  id int not null primary key,
  name varchar(255) not null,
  slug varchar(128) not null unique,
  created_at timestamp not null,
  updated_at timestamp not null
);

CREATE INDEX ON collections (updated_at);
