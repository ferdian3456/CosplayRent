package order

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/order"
	"database/sql"
	"errors"
)

type OrderRepositoryImpl struct{}

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImpl{}
}

func (repository *OrderRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, order domain.Order) {
	query := "INSERT INTO orders (id,user_id,costume_id,total,created_at) VALUES ($1,$2,$3,$4,$5)"
	_, err := tx.ExecContext(ctx, query, order.Id, order.User_id, order.Costume_id, order.Total, order.Created_at)
	helper.PanicIfError(err)
}

func (repository *OrderRepositoryImpl) FindByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.OrderResponse, error) {
	query := "SELECT id,user_id,costume_id,shipping_id,total,status_payment,is_canceled,created_at FROM orders where user_id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	orders := []order.OrderResponse{}
	for rows.Next() {
		order := order.OrderResponse{}
		err = rows.Scan(&order.User_id, &order.Costume_id, &order.Created_at)
		helper.PanicIfError(err)
		orders = append(orders, order)
		hasData = true
	}
	if hasData == false {
		return orders, errors.New("order not found")
	}

	return orders, nil
}
