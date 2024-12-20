package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"database/sql"
	"errors"

	"github.com/rs/zerolog"
)

type TopUpOrderRepository struct {
	Log *zerolog.Logger
}

func NewTopUpOrderRepository(zerolog *zerolog.Logger) *TopUpOrderRepository {
	return &TopUpOrderRepository{
		Log: zerolog,
	}
}

func (repository *TopUpOrderRepository) CreateTopUpOrder(ctx context.Context, tx *sql.Tx, topuporder domain.TopUpOrder) {
	query := "INSERT INTO topup_orders (id,user_id,topup_amount,status_payment,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := tx.ExecContext(ctx, query, topuporder.Id, topuporder.User_id, topuporder.TopUp_amount, topuporder.Status_payment, topuporder.Created_at, topuporder.Updated_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *TopUpOrderRepository) FindUserIdByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	query := "SELECT user_id FROM topup_orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var userid string
	if row.Next() {
		err = row.Scan(&userid)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return userid, nil
	} else {
		return userid, errors.New("user not found")
	}
}

func (repository TopUpOrderRepository) Update(ctx context.Context, tx *sql.Tx, midtrans domain.Midtrans) {
	query := "UPDATE topup_orders SET status_payment=true, updated_at=$1 WHERE id=$2"
	_, err := tx.ExecContext(ctx, query, midtrans.Updated_at, midtrans.Order_id)

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}
