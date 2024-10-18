package costume

import (
	"context"
	"cosplayrent/model/web/costume"
)

type CostumeService interface {
	Create(ctx context.Context, request costume.CostumeCreateRequest)
	FindById(ctx context.Context, id int) costume.CostumeResponse
	FindAll(ctx context.Context) []costume.CostumeResponse
	Update(ctx context.Context, request costume.CostumeUpdateRequest)
	Delete(ctx context.Context, id int)
}
