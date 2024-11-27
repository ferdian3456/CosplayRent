package midtrans

import (
	"context"
	"cosplayrent/helper"
	"database/sql"
	"log"
	"time"
)

type MidtransRepositoryImpl struct{}

func NewMidtransRepository() MidtransRepository {
	return &MidtransRepositoryImpl{}
}

func (repository *MidtransRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, orderid string, time *time.Time) {
	query := "UPDATE orders SET status_payment=true, updated_at=$1  WHERE id=$2"
	log.Println(query)
	_, err := tx.ExecContext(ctx, query, time, orderid)
	helper.PanicIfError(err)
}
