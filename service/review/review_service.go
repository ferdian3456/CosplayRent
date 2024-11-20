package review

import (
	"context"
	"cosplayrent/model/web/review"
)

type ReviewService interface {
	Create(ctx context.Context, request review.ReviewCreateRequest)
	FindByCostumeId(ctx context.Context, id int) []review.ReviewResponse
	FindUserReview(ctx context.Context, uuid string) []review.OwnReviewResponse
	FindUserReviewByReviewID(ctx context.Context, uuid string, reviewid int) review.OwnReviewByReviewID
	Update(ctx context.Context, request review.ReviewUpdateRequest, uuid string)
	DeleteUserReviewByReviewID(ctx context.Context, uuid string, reviewid int)
}
