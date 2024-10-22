BEGIN;

CREATE TABLE IF NOT EXISTS quizzes (
  id SERIAL PRIMARY KEY,
  requirement_id INT REFERENCES records(id),
  user_id INT REFERENCES users(id),
  name VARCHAR(255),
  description TEXT,
  threshold DECIMAL(10, 2),
  created_by_id INT REFERENCES users(id),
  updated_by_id INT REFERENCES users(id),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;