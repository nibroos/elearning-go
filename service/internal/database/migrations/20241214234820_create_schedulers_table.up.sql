BEGIN;

CREATE TABLE IF NOT EXISTS schedulers (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  cron TEXT NOT NULL,
  payload JSONB,
  status TEXT NOT NULL,
  entry_id INT,
  start_at timestamp with time zone,
  end_at timestamp with time zone,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;