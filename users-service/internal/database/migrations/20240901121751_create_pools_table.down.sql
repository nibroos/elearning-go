CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  tb1_id INT NOT NULL,
  tb2_id INT NOT NULL,
  created_by_id INT REFERENCES users(id),
  updated_by_id INT REFERENCES users(id),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT
);