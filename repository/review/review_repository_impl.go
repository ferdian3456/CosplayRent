package review

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/review"
	"database/sql"
	"errors"
)

type ReviewRepositoryImpl struct{}

func NewReviewRepository() ReviewRepository {
	return &ReviewRepositoryImpl{}
}

func (repository *ReviewRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, review domain.Review) {
	query := "INSERT INTO reviews (user_id,costume_id,description,rating,created_at) VALUES ($1,$2,$3,$4,$5)"
	_, err := tx.ExecContext(ctx, query, review.User_id, review.Costume_id, review.Description, review.Rating, review.Created_at)
	helper.PanicIfError(err)
}

func (repository *ReviewRepositoryImpl) FindByCostumeId(ctx context.Context, tx *sql.Tx, id int) ([]review.ReviewResponse, error) {
	query := "SELECT user_id,costume_id,description,rating,created_at FROM reviews where costume_id=$1"
	rows, err := tx.QueryContext(ctx, query, id)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	reviews := []review.ReviewResponse{}
	for rows.Next() {
		review := review.ReviewResponse{}
		err = rows.Scan(&review.User_id, &review.Costume_id, &review.Description, &review.Rating, &review.Created_at)
		helper.PanicIfError(err)
		reviews = append(reviews, review)
		hasData = true
	}

	if hasData == false {
		return reviews, errors.New("review not found")
	}

	return reviews, nil
}
