package review

import (
	"context"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/review"
	"database/sql"
)

type ReviewRepository interface {
	Create(ctx context.Context, tx *sql.Tx, review domain.Review)
	FindByCostumeId(ctx context.Context, tx *sql.Tx, id int) ([]review.ReviewResponse, error)
	FindUserReview(ctx context.Context, tx *sql.Tx, uuid string) ([]review.OwnReviewResponse, error)
	Update(ctx context.Context, tx *sql.Tx, review review.ReviewUpdateRequest, uuid string)
	FindUserReviewByReviewID(ctx context.Context, tx *sql.Tx, uuid string, reviewid int) (review.OwnReviewByReviewID, error)
	DeleteUserReviewByReviewID(ctx context.Context, tx *sql.Tx, uuid string, reviewid int)
}
