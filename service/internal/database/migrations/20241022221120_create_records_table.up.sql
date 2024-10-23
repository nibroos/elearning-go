BEGIN;

CREATE TABLE IF NOT EXISTS records (
  id SERIAL PRIMARY KEY,
  education_id INT REFERENCES educations(id),
  user_id INT REFERENCES users(id),
  time_spent INT,
  last_seen timestamp with time zone,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;