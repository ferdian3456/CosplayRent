package midtrans

import (
  "context"
  "database/sql"
)

type MidtransRepository interface {
	Update(ctx context.Context, tx *sql.Tx, status string)
}
