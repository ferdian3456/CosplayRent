package repository

import (
	"context"
	"cosplayrent/internal/helper"
	"database/sql"
	"github.com/rs/zerolog"
	"log"
	"time"
)

type MidtransRepository struct {
	Log *zerolog.Logger
}

func NewMidtransRepository(zerolog *zerolog.Logger) *MidtransRepository {
	return &MidtransRepository{
		Log: zerolog,
	}
}

func (repository *MidtransRepository) Update(ctx context.Context, tx *sql.Tx, orderid string, time *time.Time) {
	query := "UPDATE orders SET status_payment=true, updated_at=$1  WHERE id=$2"
	log.Println(query)
	_, err := tx.ExecContext(ctx, query, time, orderid)
	helper.PanicIfError(err)
}
