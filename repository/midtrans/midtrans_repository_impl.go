package midtrans

import (
	"context"
	"cosplayrent/helper"
	"database/sql"
	"log"
)

type MidtransRepositoryImpl struct{}

func NewMidtransRepository() MidtransRepository {
	return &MidtransRepositoryImpl{}
}

func (repository *MidtransRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, orderid string) {
	query := "UPDATE orders SET status_payment=true  WHERE id=$1"
	log.Println(query)
	_, err := tx.ExecContext(ctx, query, orderid)
	helper.PanicIfError(err)
}
