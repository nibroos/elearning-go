BEGIN;

CREATE TABLE IF NOT EXISTS contacts (
  id SERIAL PRIMARY KEY,
  type_contact_id INT REFERENCES mix_values(id),
  user_id INT REFERENCES users(id),
  ref_num INT,
  stasus INT,
  options_json JSONB,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;