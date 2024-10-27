BEGIN;

CREATE TABLE IF NOT EXISTS classes (
  id SERIAL PRIMARY KEY,
  subscribe_id INT REFERENCES subscribes(id),
  incharge_id INT REFERENCES users(id),
  name VARCHAR(255),
  description TEXT,
  banner_url TEXT,
  logo_url TEXT,
  video_url TEXT,
  created_by_id INT REFERENCES users(id),
  updated_by_id INT REFERENCES users(id),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;