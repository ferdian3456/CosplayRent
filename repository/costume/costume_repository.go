package costume

import (
	"context"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/costume"
	"database/sql"
)

type CostumeRepository interface {
	Create(ctx context.Context, tx *sql.Tx, costume domain.Costume)
	FindById(ctx context.Context, tx *sql.Tx, id int) (costume.CostumeResponse, error)
	FindByName(ctx context.Context, tx *sql.Tx, name string) ([]costume.CostumeResponse, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]costume.CostumeResponse, error)
	Update(ctx context.Context, tx *sql.Tx, costume costume.CostumeUpdateRequest)
	Delete(ctx context.Context, tx *sql.Tx, id int)
	FindByUserUUID(ctx context.Context, tx *sql.Tx, userUUID string) ([]costume.CostumeResponse, error)
	FindSellerCostumeByCostumeID(ctx context.Context, tx *sql.Tx, userUUID string, costumeID int) (costume.CostumeResponse, error)
}
