package midtrans

import (
	"context"
	"database/sql"
	"time"
)

type MidtransRepository interface {
	Update(ctx context.Context, tx *sql.Tx, orderid string, time *time.Time)
}
