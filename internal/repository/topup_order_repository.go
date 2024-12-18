package repository

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web/topup_order"
	"cosplayrent/internal/model/web/user"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
	"log"
	"time"
)

type TopUpOrderRepository struct {
	Log *zerolog.Logger
}

func NewTopUpOrderRepository(zerolog *zerolog.Logger) *TopUpOrderRepository {
	return &TopUpOrderRepository{
		Log: zerolog,
	}
}

func (repository *TopUpOrderRepository) CreateTopUpOrder(ctx context.Context, tx *sql.Tx, orderid string, uuid string, emoney user.TopUpEmoney, time *time.Time) {
	log.Printf("User with uuid: %s enter User Repository: CreateTopUpOrder", uuid)
	query := "INSERT INTO topup_orders (id,user_id,topup_amount,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)"
	_, err := tx.ExecContext(ctx, query, orderid, uuid, emoney.Emoney_amont, time, time)
	helper.PanicIfError(err)
}

func (t TopUpOrderRepository) FindUserIdByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (userId string, err error) {
	log.Println("Midtrans CallBack enter TopUpOrderRepository: FindUserIdByOrderId")

	query := "SELECT user_id FROM topup_orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderId)
	helper.PanicIfError(err)

	defer row.Close()

	var userid string

	if row.Next() {
		err = row.Scan(&userid)
		helper.PanicIfError(err)
		return userid, nil
	} else {
		return "", errors.New("topuporder is not found")
	}
}

func (t TopUpOrderRepository) UpdateTopUpOrder(ctx context.Context, tx *sql.Tx, topuporderid string, time *time.Time) {
	log.Println("Midtrans callback enter TopUpOrderRepository: UpdateTopUpOrder")

	query := "UPDATE topup_orders SET status_payment=true, updated_at=$1 WHERE id=$2"
	_, err := tx.ExecContext(ctx, query, time, topuporderid)

	helper.PanicIfError(err)
}

func (repository *TopUpOrderRepository) FindTopUpOrderHistoryByUserId(ctx context.Context, tx *sql.Tx, userid string) ([]user.UserEmoneyResponse, error) {
	log.Println("User with uuid: %s enter User Repository: CreateTopUpOrder", userid)

	query := "SELECT topup_amount,updated_at FROM topup_orders WHERE user_id=$1 AND status_payment=true"
	rows, err := tx.QueryContext(ctx, query, userid)
	helper.PanicIfError(err)

	defer rows.Close()

	topUpOrders := []user.UserEmoneyResponse{}
	var updatedAt time.Time
	var hasData bool = false

	for rows.Next() {
		topUpOrder := user.UserEmoneyResponse{}
		err = rows.Scan(&topUpOrder.Emoney_amont, &updatedAt)
		topUpOrder.Emoney_updated_at = updatedAt.Format("2006-01-02 15:04:05")
		helper.PanicIfError(err)
		topUpOrders = append(topUpOrders, topUpOrder)
		hasData = true
	}

	if hasData == false {
		return topUpOrders, errors.New("topup_order is not found")
	}

	return topUpOrders, nil
}

func (repository *TopUpOrderRepository) CheckTopUpOrderByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (topup_order.TopupOrderResponse, error) {
	query := "SELECT status_payment FROM topup_orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderId)
	helper.PanicIfError(err)

	defer row.Close()

	orderResult := topup_order.TopupOrderResponse{}

	if row.Next() {
		err = row.Scan(&orderResult.TopUpStatusPayment)
		helper.PanicIfError(err)
		return orderResult, nil
	} else {
		return topup_order.TopupOrderResponse{}, errors.New("topuporder is not found")
	}
}
