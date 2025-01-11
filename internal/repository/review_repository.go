package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/review"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type ReviewRepository struct {
	Log *zerolog.Logger
}

func NewReviewRepository(zerolog *zerolog.Logger) *ReviewRepository {
	return &ReviewRepository{
		Log: zerolog,
	}
}

func (repository *ReviewRepository) Create(ctx context.Context, tx *sql.Tx, review domain.Review) {
	query := "INSERT INTO reviews (customer_id,costume_id,order_id,description,review_picture,rating,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"
	_, err := tx.ExecContext(ctx, query, review.User_id, review.Costume_id, review.Order_id, review.Description, review.Review_picture, review.Rating, review.Created_at, review.Updated_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *ReviewRepository) CheckReviewId(ctx context.Context, tx *sql.Tx, reviewid int) error {
	query := "SELECT id FROM reviews WHERE id=$1"
	rows, err := tx.QueryContext(ctx, query, reviewid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	if rows.Next() {
		return nil
	} else {
		return errors.New("review not found")
	}
}

func (repository *ReviewRepository) Update(ctx context.Context, tx *sql.Tx, review domain.Review) {
	query := "UPDATE reviews SET "
	args := []interface{}{}
	argCounter := 1

	if review.Review_picture != "" {
		query += fmt.Sprintf("review_picture = $%d, ", argCounter)
		args = append(args, review.Review_picture)
		argCounter++
	}
	if review.Description != "" {
		query += fmt.Sprintf("description = $%d, ", argCounter)
		args = append(args, review.Description)
		argCounter++
	}
	if review.Rating != 0 {
		query += fmt.Sprintf("rating = $%d, ", argCounter)
		args = append(args, review.Rating)
		argCounter++
	}

	query += fmt.Sprintf("updated_at = $%d ", argCounter)
	args = append(args, review.Updated_at)
	argCounter++

	query += fmt.Sprintf("WHERE id = $%d", argCounter)
	args = append(args, review.Id)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *ReviewRepository) FindByCostumeId(ctx context.Context, tx *sql.Tx, id int) ([]review.ReviewResponse, error) {
	query := "SELECT customer_id,costume_id,review_picture,description,rating,created_at,updated_at FROM reviews where costume_id=$1"
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	hasData := false

	defer rows.Close()

	reviews := []review.ReviewResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		review := review.ReviewResponse{}
		err = rows.Scan(&review.User_id, &review.Costume_id, &review.Review_picture, &review.Description, &review.Rating, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		review.Created_at = createdAt.Format("2006-01-02 15:04:05")
		review.Updated_at = createdAt.Format("2006-01-02 15:04:05")
		reviews = append(reviews, review)
		hasData = true
	}

	if hasData == false {
		return reviews, errors.New("review not found")
	}

	return reviews, nil
}

func (repository *ReviewRepository) FindReviewPictureById(ctx context.Context, tx *sql.Tx, orderid string) *string {
	query := "SELECT review_picture FROM reviews where order_id=$1"
	rows, err := tx.QueryContext(ctx, query, orderid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	var review_picture *string
	if rows.Next() {
		err = rows.Scan(&review_picture)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return review_picture
	} else {
		return nil
	}

}
