CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  created_by_id INT REFERENCES users(id),
  updated_by_id INT REFERENCES users(id),
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone DEFAULT
);

-- Adding indexes for foreign key columns to improve query performance
CREATE INDEX idx_roles_created_by_id ON roles(created_by_id);

CREATE INDEX idx_roles_updated_by_id ON roles(updated_by_id);