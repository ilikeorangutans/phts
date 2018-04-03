CREATE TABLE share_rendition_configurations (
  share_id BIGINT NOT NULL REFERENCES shares(id) ON DELETE CASCADE,
  rendition_configuration_id BIGINT NOT NULL REFERENCES rendition_configurations(id) ON DELETE CASCADE,
  PRIMARY KEY(share_id, rendition_configuration_id),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
