create table service_users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(128) NOT NULL UNIQUE,
  password VARCHAR(60) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  last_login TIMESTAMP
);

create index on service_users (email);
