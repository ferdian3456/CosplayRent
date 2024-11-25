package review

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/review"
	"database/sql"
	"errors"
	"log"
	"time"
)

type ReviewRepositoryImpl struct {
}

func NewReviewRepository() ReviewRepository {
	return &ReviewRepositoryImpl{}
}

func (repository *ReviewRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, review domain.Review) {
	log.Printf("User with uuid: %s enter Review Controller: Create", review.User_id)
	query := "INSERT INTO reviews (user_id,costume_id,description,review_picture,rating,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.ExecContext(ctx, query, review.User_id, review.Costume_id, review.Description, review.Review_picture, review.Rating, review.Created_at, review.Created_at)
	helper.PanicIfError(err)
}

func (repository *ReviewRepositoryImpl) FindByCostumeId(ctx context.Context, tx *sql.Tx, id int) ([]review.ReviewResponse, error) {
	query := "SELECT user_id,costume_id,description,rating,created_at,updated_at FROM reviews where costume_id=$1"
	rows, err := tx.QueryContext(ctx, query, id)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	reviews := []review.ReviewResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		review := review.ReviewResponse{}
		err = rows.Scan(&review.User_id, &review.Costume_id, &review.Description, &review.Rating, &createdAt, &updatedAt)
		helper.PanicIfError(err)
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

func (repository *ReviewRepositoryImpl) FindUserReview(ctx context.Context, tx *sql.Tx, uuid string) ([]review.OwnReviewResponse, error) {
	log.Printf("User with uuid: %s enter Review Repository: FindUserReview", uuid)
	query := "SELECT id,user_id,costume_id, description, rating, review_picture, created_at, updated_at FROM reviews WHERE user_id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)

	defer rows.Close()
	hasData := false

	reviews := []review.OwnReviewResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		review := review.OwnReviewResponse{}
		err = rows.Scan(&review.Id, &review.Userid, &review.Costume_Id, &review.Description, &review.Rating, &review.Review_picture, &createdAt, &updatedAt)
		helper.PanicIfError(err)
		review.Created_at = createdAt.Format("2006-01-02 15:04:05")
		review.Updated_at = updatedAt.Format("2026-01-02 15:04:05")
		reviews = append(reviews, review)
		hasData = true
	}
	if hasData == false {
		return reviews, errors.New("review not found")
	}
	return reviews, nil
}

func (repository *ReviewRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, review review.ReviewUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter Review Repository: Update", uuid)
	if review.Review_picture == nil {
		query := "UPDATE reviews SET description=$2,rating=$3,updated_at=$4  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, review.ReviewId, review.Description, review.Rating, review.Updated_at)
		helper.PanicIfError(err)
	} else {
		query := "UPDATE reviews SET review_picture=$2,description=$3,rating=$4,updated_at=$5  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, review.ReviewId, review.Review_picture, review.Description, review.Rating, review.Updated_at)
		helper.PanicIfError(err)
	}
}

func (repository *ReviewRepositoryImpl) FindUserReviewByReviewID(ctx context.Context, tx *sql.Tx, uuid string, reviewid int) (review.OwnReviewByReviewID, error) {
	log.Printf("User with uuid: %s enter Review Repository: FindUserReviewByReviewID", uuid)
	query := "SELECT user_id, costume_id, description,review_picture,rating, created_at, updated_at FROM reviews WHERE id=$1"
	rows, err := tx.QueryContext(ctx, query, reviewid)
	helper.PanicIfError(err)
	defer rows.Close()

	review := review.OwnReviewByReviewID{}
	var createdAt time.Time
	var updatedAt time.Time
	if rows.Next() {
		err := rows.Scan(&review.User_id, &review.Costume_id, &review.Description, &review.Review_picture, &review.Rating, &createdAt, &updatedAt)
		helper.PanicIfError(err)
		review.Created_at = createdAt.Format("2006-01-02 15:04:05")
		review.Updated_at = createdAt.Format("2006-01-02 15:04:05")
		return review, nil
	} else {
		return review, errors.New("review not found")
	}

	return review, nil
}

func (repository *ReviewRepositoryImpl) DeleteUserReviewByReviewID(ctx context.Context, tx *sql.Tx, uuid string, reviewid int) {
	log.Printf("User with uuid: %s enter Review Repository: DeleteUserReviewByReviewID", uuid)
	query := "DELETE FROM reviews WHERE id=$1"
	_, err := tx.ExecContext(ctx, query, reviewid)
	helper.PanicIfError(err)
}
