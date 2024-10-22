BEGIN;

CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  type_payment_id INT REFERENCES mix_values(id),
  user_id INT REFERENCES users(id),
  petugas_id INT REFERENCES users(id),
  subscribe_id INT REFERENCES subscribes(id),
  amount DECIMAL(20, 5),
  invoice_url TEXT,
  paid_at timestamp with time zone,
  created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone,
  deleted_at timestamp with time zone
);

COMMIT;