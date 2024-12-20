package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"database/sql"
	"errors"

	"github.com/rs/zerolog"
)

type OrderRepository struct {
	Log *zerolog.Logger
}

func NewOrderRepository(zerolog *zerolog.Logger) *OrderRepository {
	return &OrderRepository{
		Log: zerolog,
	}
}

func (repository *OrderRepository) Create(ctx context.Context, tx *sql.Tx, order domain.Order) {
	query := "INSERT INTO orders (id,user_id,seller_id,costume_id,total,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.ExecContext(ctx, query, order.Id, order.User_id, order.Seller_id, order.Costume_id, order.Total, order.Created_at, order.Created_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *OrderRepository) FindBuyerIdByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	query := "SELECT user_id from orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var buyerid string

	if row.Next() {
		err = row.Scan(&buyerid)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return buyerid, nil
	} else {
		return buyerid, errors.New("buyer not found")
	}
}

func (repository *OrderRepository) FindSellerIdByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	query := "SELECT seller_id from orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var sellerid string

	if row.Next() {
		err = row.Scan(&sellerid)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return sellerid, nil
	} else {
		return sellerid, errors.New("buyer not found")
	}
}

func (repository *OrderRepository) Update(ctx context.Context, tx *sql.Tx, midtrans domain.Midtrans) {
	query := "UPDATE orders SET status_payment=true, updated_at=$1  WHERE id=$2"
	_, err := tx.ExecContext(ctx, query, midtrans.Updated_at, midtrans.Order_id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}
