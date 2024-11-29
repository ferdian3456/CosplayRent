CREATE TABLE IF NOT EXISTS topup_orders(
  id uuid PRIMARY KEY,
  user_id UUID NOT NULL,
  topup_amount decimal(10,2) NOT NULL,
  status_payment bool default false,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);