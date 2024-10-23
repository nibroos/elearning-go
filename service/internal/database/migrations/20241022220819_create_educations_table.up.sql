BEGIN;

CREATE TABLE IF NOT EXISTS educations (
  id SERIAL PRIMARY KEY,
  module_id INT REFERENCES modules(id),
  no_urut INT,
  name VARCHAR(255),
  description TEXT,
  text_materi TEXT,
  video_url TEXT,
  thumbnail_url TEXT,
  attachment_urls JSONB,
  created_by_id INT REFERENCES users(id),
  updated_by_id INT REFERENCES users(id),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;