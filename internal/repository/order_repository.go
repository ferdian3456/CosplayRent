package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/order"
	"database/sql"
	"errors"
	"time"

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

func (repository *OrderRepository) Create(ctx context.Context, tx *sql.Tx, userRequest order.DirectlyOrderToMidtrans) {
	query := "INSERT INTO orders (id,user_id,seller_id,costume_id,total,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.ExecContext(ctx, query, userRequest.Id, userRequest.Costumer_id, userRequest.Seller_id, userRequest.Costume_id, userRequest.TotalAmount, userRequest.Created_at, userRequest.Updated_at)
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

func (repository *OrderRepository) FindOrderDetailByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (order.OrderResponse, error) {
	query := "SELECT id, user_id, seller_id,description,costume_id, total, status_payment, status_shipping, is_cancelled, created_at, updated_at FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	order := order.OrderResponse{}
	var createdAt time.Time
	var updatedAt time.Time

	if row.Next() {
		err = row.Scan(&order.Id, &order.User_id, &order.Seller_id, &order.Description, &order.Costume_id, &order.Total, &order.Status_payment, &order.Status_shipping, &order.Is_canceled, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		order.Created_at = createdAt.Format("2006-01-02 15:04:05")
		order.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		return order, nil
	} else {
		return order, errors.New("order not found")
	}
}

func (repository *OrderRepository) CheckStatusPayment(ctx context.Context, tx *sql.Tx, orderid string) (*bool, error) {
	query := "SELECT status_payment FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var statusPayment *bool

	if row.Next() {
		err = row.Scan(&statusPayment)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return statusPayment, nil
	} else {
		return statusPayment, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindUserAndCostumeById(ctx context.Context, tx *sql.Tx, orderid string) (domain.Order, error) {
	query := "SELECT user_id,costume_id FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	order := domain.Order{}

	if row.Next() {
		err = row.Scan(&order.User_id, &order.Costume_id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return order, nil
	} else {
		return order, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindSellerAndCostumeById(ctx context.Context, tx *sql.Tx, orderid string) (domain.Order, error) {
	query := "SELECT seller_id,costume_id,description FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	order := domain.Order{}

	if row.Next() {
		err = row.Scan(&order.Seller_id, &order.Costume_id, &order.Description)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return order, nil
	} else {
		return order, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindOrderBySellerId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.AllSellerOrderResponse, error) {
	query := "SELECT id,costume_id,total,status,updated_at FROM orders where seller_id=$1 AND status_payment=true"
	rows, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
	hasData := false

	defer rows.Close()

	orders := []order.AllSellerOrderResponse{}
	for rows.Next() {
		order := order.AllSellerOrderResponse{}
		err = rows.Scan(&order.Id, &order.Costume_id, &order.Total, &order.Status, &order.Updated_at)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		orders = append(orders, order)
		hasData = true
	}
	if hasData == false {
		return orders, errors.New("order not found")
	}

	return orders, nil
}

func (repository *OrderRepository) FindOrderByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.AllUserOrderResponse, error) {
	query := "SELECT id,costume_id,total,status,updated_at FROM orders where user_id=$1 AND status_payment=true"
	rows, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
	hasData := false

	defer rows.Close()

	orders := []order.AllUserOrderResponse{}
	for rows.Next() {
		order := order.AllUserOrderResponse{}
		err = rows.Scan(&order.Id, &order.Costume_id, &order.Total, &order.Status, &order.Updated_at)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		orders = append(orders, order)
		hasData = true
	}
	if hasData == false {
		return orders, errors.New("order not found")
	}

	return orders, nil
}
