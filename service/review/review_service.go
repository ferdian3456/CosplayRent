package review

import (
	"context"
	"cosplayrent/model/web/review"
)

type ReviewService interface {
	Create(ctx context.Context, request review.ReviewCreateRequest)
	FindByCostumeId(ctx context.Context, id int) []review.ReviewResponse
}
