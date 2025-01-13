package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/model/web/review"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
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

func (repository *OrderRepository) Create(ctx context.Context, tx *sql.Tx, userRequest domain.Order) {
	query := "INSERT INTO orders (id,customer_id,seller_id,costume_id,total,shipment_origin,shipment_destination,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	_, err := tx.ExecContext(ctx, query, userRequest.Id, userRequest.Costumer_id, userRequest.Seller_id, userRequest.Costume_id, userRequest.Total_amount, userRequest.Shipment_origin, userRequest.Shipment_destination, userRequest.Created_at, userRequest.Updated_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *OrderRepository) FindBuyerIdByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	query := "SELECT customer_id from orders WHERE id=$1"
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

func (repository *OrderRepository) CheckOrderAndCostumeId(ctx context.Context, tx *sql.Tx, orderid string, costumeid int) error {
	query := "SELECT id,costume_id from orders WHERE id=$1 AND costume_id=$2"
	row, err := tx.QueryContext(ctx, query, orderid, costumeid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	if row.Next() {
		return nil
	} else {
		return errors.New("order not found")
	}
}

func (repository *OrderRepository) Update(ctx context.Context, tx *sql.Tx, midtrans domain.Midtrans) {
	query := "UPDATE payments SET status='Paid', updated_at=$1  WHERE order_id=$2"
	_, err := tx.ExecContext(ctx, query, midtrans.Updated_at, midtrans.Order_id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *OrderRepository) CreatePayment(ctx context.Context, tx *sql.Tx, payment domain.Payments) {
	query := "INSERT INTO payments (order_id,customer_id,seller_id,status,amount,method,midtrans_redirect_url,midtrans_url_expired_time,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	_, err := tx.ExecContext(ctx, query, payment.Order_id, payment.Customer_id, payment.Seller_id, payment.Status, payment.Amount, payment.Payment_method, payment.Midtrans_redirect_url, payment.Midtrans_url_expired_time, payment.Created_at, payment.Updated_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *OrderRepository) FindPaymentMethodByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	query := "SELECT method FROM payments WHERE order_id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var paymentMethod string

	if row.Next() {
		err = row.Scan(&paymentMethod)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return paymentMethod, nil
	} else {
		return paymentMethod, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindPaymentInfoByPaymentId(ctx context.Context, tx *sql.Tx, paymentid int, customerid string) (domain.Payments, error) {
	query := "SELECT amount,status,midtrans_redirect_url,midtrans_url_expired_time,created_at FROM payments WHERE id=$1 AND customer_id=$2"
	row, err := tx.QueryContext(ctx, query, paymentid, customerid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	payment := domain.Payments{}

	if row.Next() {
		err = row.Scan(&payment.Amount, &payment.Status, &payment.Midtrans_redirect_url, &payment.Midtrans_url_expired_time, &payment.Created_at)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}

		return payment, nil
	} else {
		return payment, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindPaymentInfoByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (domain.Payments, error) {
	query := "SELECT id,status,midtrans_url_expired_time,created_at FROM payments WHERE order_id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	payment := domain.Payments{}

	if row.Next() {
		err = row.Scan(&payment.Id, &payment.Status, &payment.Midtrans_url_expired_time, &payment.Created_at)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}

		return payment, nil
	} else {
		return payment, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindOrderDetailByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (order.OrderResponse, error) {
	query := "SELECT id, customer_id, seller_id,description,costume_id, total, status_payment, status_shipping, is_cancelled, created_at, updated_at FROM orders WHERE id=$1"
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
		err = row.Scan(&order.Id, &order.Customer_id, &order.Seller_id, &order.Description, &order.Costume_id, &order.Total, &order.Status_payment, &order.Status_shipping, &order.Is_canceled, &createdAt, &updatedAt)
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

func (repository *OrderRepository) CheckStatusPayment(ctx context.Context, tx *sql.Tx, orderid string) (string, error) {
	query := "SELECT status FROM payments WHERE order_id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var statusPayment *string

	if row.Next() {
		err = row.Scan(&statusPayment)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return *statusPayment, nil
	} else {
		return *statusPayment, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindUserAndCostumeById(ctx context.Context, tx *sql.Tx, orderid string) (domain.Order, error) {
	query := "SELECT shipment_destination,customer_id,costume_id FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	order := domain.Order{}

	if row.Next() {
		err = row.Scan(&order.Shipment_destination, &order.Costumer_id, &order.Costume_id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return order, nil
	} else {
		return order, errors.New("order not found")
	}
}

func (repository *OrderRepository) GetEventDetail(ctx context.Context, tx *sql.Tx, orderid string) domain.OrderEvents {
	query := "SELECT notes,shipment_receipt_user_id FROM order_events WHERE order_id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	order := domain.OrderEvents{}

	if row.Next() {
		err = row.Scan(&order.Notes, &order.Shipment_receipt_user_id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return order
	} else {
		return order
	}
}

func (repository *OrderRepository) CreateOrderEvents(ctx context.Context, tx *sql.Tx, events domain.OrderEvents) {
	query := "INSERT INTO order_events (user_id,order_id,status,notes,shipment_receipt_user_id,created_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := tx.ExecContext(ctx, query, events.User_id, events.Order_id, events.Status, events.Notes, events.Shipment_receipt_user_id, events.Created_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *OrderRepository) FindSellerAndCostumeById(ctx context.Context, tx *sql.Tx, orderid string) (domain.Order, error) {
	query := "SELECT id,seller_id,costume_id FROM orders WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	order := domain.Order{}

	if row.Next() {
		err = row.Scan(&order.Id, &order.Seller_id, &order.Costume_id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return order, nil
	} else {
		return order, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindEventInfoById(ctx context.Context, tx *sql.Tx, orderid string) (domain.OrderEvents, error) {
	query := "SELECT notes,shipment_receipt_user_id FROM order_events WHERE order_id=$1 ORDER BY created_at DESC LIMIT 1"
	row, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	orderEvent := domain.OrderEvents{}

	if row.Next() {
		err = row.Scan(&orderEvent.Notes, &orderEvent.Shipment_receipt_user_id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}

		return orderEvent, nil
	} else {
		return orderEvent, errors.New("order not found")
	}
}

func (repository *OrderRepository) FindOrderBySellerId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.AllSellerOrderResponse, error) {
	query := `
    SELECT DISTINCT ON (o.id) 
    o.id AS order_id, 
    o.costume_id, 
    o.total, 
    e.status, 
    o.updated_at
	FROM 
		orders o
	JOIN 
		order_events e 
	ON 
		o.id = e.order_id
	WHERE 
		o.seller_id = $1
	ORDER BY 
		o.id, e.created_at DESC;
	`

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
	query := `
    SELECT DISTINCT ON (o.id) 
    o.id AS order_id, 
    o.costume_id, 
    o.total, 
    e.status, 
    o.updated_at
	FROM 
		orders o
	JOIN 
		order_events e 
	ON 
		o.id = e.order_id
	WHERE 
		o.customer_id = $1
	ORDER BY 
		o.id, e.created_at DESC;
	`
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

func (repository *OrderRepository) FindOrderInfoByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]review.UserReviewResponse, error) {
	query := "SELECT id,seller_id,costume_id FROM orders where customer_id=$1 AND status_payment=true"
	rows, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
	hasData := false

	defer rows.Close()

	reviews := []review.UserReviewResponse{}
	for rows.Next() {
		review := review.UserReviewResponse{}
		err = rows.Scan(&review.Order_id, &review.Seller_id, &review.Custome_Id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		reviews = append(reviews, review)
		hasData = true
	}
	if hasData == false {
		return reviews, errors.New("order not found")
	}

	return reviews, nil
}

func (repository *OrderRepository) FindListOrderByCostumeId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.PaymentTransationForOrderResponse, error) {
	query := "SELECT id,costume_id FROM orders where customer_id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
	hasData := false

	defer rows.Close()

	orders := []order.PaymentTransationForOrderResponse{}
	for rows.Next() {
		order := order.PaymentTransationForOrderResponse{}
		err = rows.Scan(&order.Order_id, &order.Costume_id)
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

func (repository *OrderRepository) CheckIfUserOrSeller(ctx context.Context, tx *sql.Tx, userid string, orderid string) error {
	query := "SELECT customer_id FROM orders where customer_id=$1 AND status_payment=true AND id=$2"
	row, err := tx.QueryContext(ctx, query, userid, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	if row.Next() {
		return nil
	} else {
		return errors.New("order not found")
	}
}

func (repository *OrderRepository) UpdateOrder(ctx context.Context, tx *sql.Tx, order domain.Order) {
	query := "UPDATE orders SET status=$1,description=$2 WHERE id=$3 "
	_, err := tx.ExecContext(ctx, query, order.Id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}
