CREATE TABLE exif (
  id SERIAL PRIMARY KEY,
  photo_id INTEGER NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
  value_type INTEGER NOT NULL,
  tag VARCHAR(128) NOT NULL,
  string VARCHAR(256) NOT NULL,
  num BIGINT,
  denom INTEGER,
  datetime TIMESTAMP,
  floating DOUBLE PRECISION
);

CREATE INDEX ON exif (photo_id, tag, datetime);
