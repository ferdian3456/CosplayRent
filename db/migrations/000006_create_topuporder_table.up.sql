CREATE TABLE IF NOT EXISTS topup_orders(
  id char(36) PRIMARY KEY,
  user_id char(36) NOT NULL,
  topup_amount decimal(10,2) NOT NULL,
  status_payment varchar(7) DEFAULT 'Pending',
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);