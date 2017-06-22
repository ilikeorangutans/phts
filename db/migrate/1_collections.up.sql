CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  handle VARCHAR(32) NOT NULL UNIQUE,
  email VARCHAR(128) NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON users (updated_at);

CREATE TABLE collections (
  id serial primary key,
  name varchar(255) not null,
  slug varchar(128) not null unique,
  created_at timestamp not null,
  updated_at timestamp not null
);

CREATE INDEX ON collections (updated_at);

CREATE TABLE photos (
  id SERIAL PRIMARY KEY,
  collection_id INTEGER NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
  description TEXT,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON photos (updated_at);
CREATE INDEX ON photos (collection_id, updated_at);

CREATE TABLE renditions (
  id SERIAL PRIMARY KEY,
  photo_id INTEGER NOT NULL REFERENCES photos(id) ON DELETE CASCADE,

  original BOOLEAN DEFAULT false NOT NULL,
  width INTEGER NOT NULL,
  height INTEGER NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON photos (updated_at);
CREATE INDEX ON renditions (photo_id, updated_at);

