BEGIN;

CREATE TABLE IF NOT EXISTS groups (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  status INT,
  options_json JSONB,
  created_by_id INT REFERENCES users(id),
  updated_by_id INT REFERENCES users(id),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone
);

COMMIT;