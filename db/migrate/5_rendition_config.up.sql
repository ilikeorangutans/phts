CREATE TABLE rendition_configurations (
  id SERIAL PRIMARY KEY,
  collection_id INTEGER REFERENCES collections(id) ON DELETE CASCADE,
  name VARCHAR(128) NOT NULL DEFAULT '',

  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  quality INTEGER NOT NULL DEFAULT 95,

  created_at TIMESTAMP NOT NULL
);

INSERT INTO rendition_configurations (name, width, height, quality, created_at) VALUES ('admin thumbnails', 345, 0, 95, NOW()), ('admin preview', 635, 0, NOW());
