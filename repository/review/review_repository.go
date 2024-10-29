package review

import (
	"context"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/review"
	"database/sql"
)

type ReviewRepository interface {
	Create(ctx context.Context, tx *sql.Tx, costume domain.Review)
	FindByCostumeId(ctx context.Context, tx *sql.Tx, id int) ([]review.ReviewResponse, error)
}
