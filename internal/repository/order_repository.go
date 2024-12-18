package repository

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/model/web/user"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
	"log"
	"time"
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
	log.Println(order.Seller_id)
	query := "INSERT INTO orders (id,user_id,seller_id,costume_id,total,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.ExecContext(ctx, query, order.Id, order.User_id, order.Seller_id, order.Costume_id, order.Total, order.Created_at, order.Created_at)
	helper.PanicIfError(err)
}

func (repository *OrderRepository) FindByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.OrderResponse, error) {
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

func (repository *OrderRepository) FindOrderBySellerId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.AllSellerOrderResponse, error) {
	query := "SELECT id,costume_id,total,status,updated_at FROM orders where seller_id=$1 AND status_payment=true"
	rows, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	orders := []order.AllSellerOrderResponse{}
	for rows.Next() {
		order := order.AllSellerOrderResponse{}
		err = rows.Scan(&order.Id, &order.Costume_id, &order.Total, &order.Status, &order.Updated_at)
		helper.PanicIfError(err)
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
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	orders := []order.AllUserOrderResponse{}
	for rows.Next() {
		order := order.AllUserOrderResponse{}
		err = rows.Scan(&order.Id, &order.Costume_id, &order.Total, &order.Status, &order.Updated_at)
		helper.PanicIfError(err)
		orders = append(orders, order)
		hasData = true
	}
	if hasData == false {
		return orders, errors.New("order not found")
	}

	return orders, nil
}

func (repository *OrderRepository) DirectlyOrderToMidtrans(ctx context.Context, tx *sql.Tx, uuid string, sendOrderToDatabase order.DirectlyOrderToMidtrans) {
	log.Printf("User with uuid: %s enter Order Repository: DirectlyOrderToMidtrans", uuid)

	query := "INSERT INTO orders (id,user_id,seller_id,costume_id,total,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.ExecContext(ctx, query, sendOrderToDatabase.Id, sendOrderToDatabase.Costumer_id, sendOrderToDatabase.Seller_id, sendOrderToDatabase.Costume_id, sendOrderToDatabase.TotalAmount, sendOrderToDatabase.Created_at, sendOrderToDatabase.Created_at)
	helper.PanicIfError(err)
}

func (repository *OrderRepository) FindBuyerIdByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	log.Println("From Midtrans Callback enter Order Repository: FindBuyerIdByOrderId")

	query := "SELECT user_id FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	helper.PanicIfError(err)

	defer row.Close()

	var UserId string
	if row.Next() {
		err = row.Scan(&UserId)
		helper.PanicIfError(err)
		return UserId, nil
	} else {
		return "", errors.New("order not found")
	}
}

func (repository *OrderRepository) FindSellerIdByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	log.Println("From Midtrans Callback enter Order Repository: FindSellerIdByOrderId")

	query := "SELECT seller_id FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	helper.PanicIfError(err)

	defer row.Close()

	var SellerId string
	if row.Next() {
		err = row.Scan(&SellerId)
		helper.PanicIfError(err)
		return SellerId, nil
	} else {
		return "", errors.New("order not found")
	}
}

func (repository *OrderRepository) FindOrderDetailByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (order.OrderResponse, error) {
	query := "SELECT id, user_id, seller_id,description,costume_id, total, status_payment, status_shipping, is_cancelled, created_at, updated_at FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	helper.PanicIfError(err)

	defer row.Close()

	order := order.OrderResponse{}
	var createdAt time.Time
	var updatedAt time.Time

	if row.Next() {
		err = row.Scan(&order.Id, &order.User_id, &order.Seller_id, &order.Description, &order.Costume_id, &order.Total, &order.Status_payment, &order.Status_shipping, &order.Is_canceled, &createdAt, &updatedAt)
		helper.PanicIfError(err)
		order.Created_at = createdAt.Format("2006-01-02 15:04:05")
		order.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
	} else {
		return order, errors.New("order not found")
	}

	return order, nil
}

func (repository *OrderRepository) FindOrderHistoryByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserEmoneyResponse, error) {
	log.Println("User with : FindSellerIdByOrderId")

	query := "SELECT total,updated_at FROM orders WHERE user_id=$1 AND status_payment=true"
	rows, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)

	defer rows.Close()

	orders := []user.UserEmoneyResponse{}
	var updatedAt time.Time
	var hasData bool = false

	for rows.Next() {
		order := user.UserEmoneyResponse{}
		err = rows.Scan(&order.Emoney_amont, &updatedAt)
		helper.PanicIfError(err)
		order.Emoney_updated_at = updatedAt.Format("2006-01-02 15:04:05")
		orders = append(orders, order)
		hasData = true
	}

	if hasData == false {
		return orders, errors.New("order not found")
	}

	return orders, nil
}

func (repository *OrderRepository) FindOrderHistoryBySellerId(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserEmoneyResponse, error) {
	log.Printf("User with uuid: %s enter Order Repository: FindOrderHistoryBySellerId", uuid)

	query := "SELECT total,updated_at FROM orders WHERE seller_id=$1 AND status_payment=true"
	rows, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)

	defer rows.Close()

	orders := []user.UserEmoneyResponse{}
	var updatedAt time.Time
	var hasData bool = false

	for rows.Next() {
		order := user.UserEmoneyResponse{}
		err = rows.Scan(&order.Emoney_amont, &updatedAt)
		helper.PanicIfError(err)
		order.Emoney_updated_at = updatedAt.Format("2006-01-02 15:04:05")
		orders = append(orders, order)
		hasData = true
	}

	if hasData == false {
		return orders, errors.New("order is not found")
	}

	return orders, nil
}

func (repository *OrderRepository) UpdateSellerOrder(ctx context.Context, tx *sql.Tx, updateRequest order.OrderUpdateRequest, sellerid string, orderid string) {
	log.Printf("User with uuid: %s enter Order Repository: UpdateSellerOrder", sellerid)

	if updateRequest.Description == "" {
		query := "UPDATE orders SET status=$1,updated_at=$2 WHERE id=$3"
		_, err := tx.ExecContext(ctx, query, updateRequest.StatusOrder, updateRequest.Updated_at, orderid)
		helper.PanicIfError(err)
	} else {
		query := "UPDATE orders SET status=$1,description=$2,updated_at=$3 WHERE id=$4"
		_, err := tx.ExecContext(ctx, query, updateRequest.StatusOrder, updateRequest.Description, updateRequest.Updated_at, orderid)
		helper.PanicIfError(err)
	}
}
