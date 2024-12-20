package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
)

type CostumeRepository struct {
	Log *zerolog.Logger
}

func NewCostumeRepository(zerolog *zerolog.Logger) *CostumeRepository {
	return &CostumeRepository{
		Log: zerolog,
	}
}

func (repository *CostumeRepository) Create(ctx context.Context, tx *sql.Tx, costume domain.Costume) {
	query := "INSERT INTO costumes (user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	_, err := tx.ExecContext(ctx, query, costume.User_id, costume.Name, costume.Description, costume.Bahan, costume.Ukuran, costume.Berat, costume.Kategori, costume.Price, costume.Picture, costume.Created_at, costume.Updated_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *CostumeRepository) Update(ctx context.Context, tx *sql.Tx, costume domain.Costume) {
	query := "UPDATE costumes SET "
	args := []interface{}{}
	argCounter := 1

	if costume.Name != "" {
		query += fmt.Sprintf("name = $%d, ", argCounter)
		args = append(args, costume.Name)
		argCounter++
	}
	if costume.Description != "" {
		query += fmt.Sprintf("description = $%d, ", argCounter)
		args = append(args, costume.Description)
		argCounter++
	}
	if costume.Bahan != "" {
		query += fmt.Sprintf("bahan = $%d, ", argCounter)
		args = append(args, costume.Bahan)
		argCounter++
	}
	if costume.Ukuran != "" {
		query += fmt.Sprintf("ukuran = $%d, ", argCounter)
		args = append(args, costume.Ukuran)
		argCounter++
	}
	if costume.Berat != 0 {
		query += fmt.Sprintf("berat = $%d, ", argCounter)
		args = append(args, costume.Berat)
		argCounter++
	}
	if costume.Kategori != "" {
		query += fmt.Sprintf("kategori = $%d, ", argCounter)
		args = append(args, costume.Kategori)
		argCounter++
	}
	if costume.Price != 0 {
		query += fmt.Sprintf("price = $%d, ", argCounter)
		args = append(args, costume.Price)
		argCounter++
	}
	if costume.Picture != "" {
		query += fmt.Sprintf("picture = $%d, ", argCounter)
		args = append(args, costume.Picture)
		argCounter++
	}

	query += fmt.Sprintf("updated_at = $%d ", argCounter)
	args = append(args, costume.Updated_at)
	argCounter++

	query += fmt.Sprintf("WHERE id = $%d", argCounter)
	args = append(args, costume.Id)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *CostumeRepository) CheckOwnership(ctx context.Context, tx *sql.Tx, userUUID string, costumeid int) error {
	query := "SELECT name FROM costumes WHERE id=$1 AND user_id=$2"
	rows, err := tx.QueryContext(ctx, query, costumeid, userUUID)

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	hasData := false

	defer rows.Close()

	for rows.Next() {
		hasData = true
	}

	if hasData == true {
		return errors.New("you are the owner of this costume")
	} else {
		return nil
	}
}

func (repository *CostumeRepository) FindSellerIdFindByCostumeID(ctx context.Context, tx *sql.Tx, costumeid int) (string, error) {
	query := "SELECT seller_id FROM costumes WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, costumeid)

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var seller_id string

	if row.Next() {
		err = row.Scan(&seller_id)
		if err != nil {
			respErr := errors.New("failed to query into database")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return seller_id, nil
	} else {
		return seller_id, errors.New("seller not found")
	}
}
