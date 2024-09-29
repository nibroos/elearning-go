BEGIN;

CREATE TABLE IF NOT EXISTS pools (
  id SERIAL PRIMARY KEY,
  group1_id INT REFERENCES groups(id),
  group2_id INT REFERENCES groups(id),
  mv1_id INT,
  mv2_id INT,
  description VARCHAR(255),
  options_json JSONB,
  created_by_id INT REFERENCES users(id),
  updated_by_id INT REFERENCES users(id),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;