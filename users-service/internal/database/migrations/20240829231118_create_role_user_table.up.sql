CREATE TABLE IF NOT EXISTS role_user (
  role_id INT REFERENCES roles(id) ON DELETE CASCADE,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, user_id)
);