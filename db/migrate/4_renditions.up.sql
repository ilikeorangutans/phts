CREATE TABLE renditions (
  id SERIAL PRIMARY KEY,
  photo_id INTEGER NOT NULL REFERENCES photos(id) ON DELETE CASCADE,

  width INTEGER NOT NULL,
  height INTEGER NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX ON photos (updated_at);
CREATE INDEX ON renditions (photo_id, updated_at);

ALTER TABLE photos ADD COLUMN description TEXT;
