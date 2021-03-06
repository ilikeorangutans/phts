CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(128) NOT NULL UNIQUE,
  password VARCHAR(60) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  last_login TIMESTAMP,
  must_change_password boolean not null default false
);

CREATE INDEX ON users (updated_at);

CREATE TABLE collections (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) not null,
  slug VARCHAR(128) not null unique,
  photo_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP not null,
  updated_at TIMESTAMP not null
);

CREATE INDEX ON collections (updated_at);

CREATE TABLE users_collections (
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  collection_id INTEGER NOT NULL REFERENCES collections(id) ON DELETE CASCADE,

  PRIMARY KEY(user_id, collection_id),

  created_at TIMESTAMP not null,
  updated_at TIMESTAMP not null
);

CREATE INDEX ON users_collections (updated_at);

CREATE TABLE photos (
  id SERIAL PRIMARY KEY,
  collection_id INTEGER NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
  description TEXT NOT NULL DEFAULT '',
  taken_at TIMESTAMP,
  filename VARCHAR(128) NOT NULL,
  rendition_count INTEGER NOT NULL DEFAULT 0,
  published BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON photos (updated_at);
CREATE INDEX ON photos (collection_id, updated_at);
CREATE INDEX ON photos (collection_id, taken_at);

CREATE TABLE rendition_configurations (
  id SERIAL PRIMARY KEY,
  collection_id INTEGER REFERENCES collections(id) ON DELETE CASCADE,
  name VARCHAR(128) NOT NULL DEFAULT '',
  private BOOLEAN DEFAULT FALSE NOT NULL,

  resize BOOLEAN DEFAULT TRUE NOT NULL,
  original BOOLEAN DEFAULT FALSE NOT NULL,
  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  quality INTEGER NOT NULL DEFAULT 95,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX ON rendition_configurations (collection_id, name);

INSERT INTO rendition_configurations (name, width, height, quality, original, resize, private, updated_at, created_at)
VALUES
  ('original', 0, 0, 85, true, false, true, NOW(), NOW()),
  ('admin thumbnails', 345, 0, 85, false, true, true, NOW(), NOW()),
  ('admin preview', 635, 0, 95, false, true, true, NOW(), NOW());

CREATE TABLE renditions (
  id SERIAL PRIMARY KEY,
  photo_id INTEGER NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
  rendition_configuration_id INTEGER REFERENCES rendition_configurations(id),

  original BOOLEAN DEFAULT false NOT NULL,
  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  format VARCHAR(32) NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX ON renditions (photo_id, rendition_configuration_id);
CREATE INDEX ON renditions (photo_id, updated_at);

CREATE TABLE exif (
  id SERIAL PRIMARY KEY,
  photo_id INTEGER NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
  value_type INTEGER NOT NULL,
  tag VARCHAR(128) NOT NULL,
  string VARCHAR(256) NOT NULL,
  num BIGINT,
  denom INTEGER,
  datetime TIMESTAMP,
  floating DOUBLE PRECISION,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON exif (photo_id, tag, datetime);

CREATE TABLE albums (
  id SERIAL PRIMARY KEY,
  name VARCHAR(256) NOT NULL,
  slug VARCHAR(128) NOT NULL,

  collection_id INTEGER REFERENCES collections(id) ON DELETE CASCADE,

  photo_count INTEGER NOT NULL DEFAULT 0,
  cover_photo_id INTEGER REFERENCES photos(id) ON DELETE SET NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE album_photos (
  photo_id INTEGER NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
  album_id INTEGER NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
  sort_order INTEGER NOT NULL DEFAULT 0,

  PRIMARY KEY(photo_id, album_id),

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON album_photos (photo_id, album_id, sort_order, updated_at);

CREATE TABLE share_sites (
  id SERIAL PRIMARY KEY,
  domain VARCHAR(128) NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX ON share_sites (domain);

CREATE TABLE shares (
  id SERIAL PRIMARY KEY,

  share_site_id INTEGER NOT NULL REFERENCES share_sites(id) ON DELETE CASCADE,

  photo_id INTEGER REFERENCES photos(id),
  collection_id INTEGER REFERENCES collections(id),
  slug VARCHAR(128) NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON shares (updated_at, share_site_id);
CREATE UNIQUE INDEX ON shares (share_site_id, slug);

