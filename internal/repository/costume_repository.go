package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/costume"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"time"
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
	row, err := tx.QueryContext(ctx, query, costumeid, userUUID)

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	if row.Next() {
		return errors.New("you are the owner of this costume")
	} else {
		return nil
	}
}

func (repository *CostumeRepository) CheckDeleteCostume(ctx context.Context, tx *sql.Tx, userUUID string, costumeid int) error {
	query := "SELECT name FROM costumes WHERE id=$1 AND user_id=$2"
	row, err := tx.QueryContext(ctx, query, costumeid, userUUID)

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	if row.Next() {
		return nil
	} else {
		return errors.New("not allowed to delete this costume")
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
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return seller_id, nil
	} else {
		return seller_id, errors.New("seller not found")
	}
}

func (repository *CostumeRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at, updated_at FROM costumes"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
	hasData := false

	defer rows.Close()

	costumes := []costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		costume.Created_at = createdAt.Format("2006-01-02 15:04:05")
		costume.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		costumes = append(costumes, costume)
		hasData = true
	}
	if hasData == false {
		return costumes, errors.New("costume not found")
	}

	return costumes, nil
}

func (repository *CostumeRepository) FindSellerCostume(ctx context.Context, tx *sql.Tx, uuid string) ([]costume.SellerCostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at, updated_at FROM costumes where user_id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
	hasData := false

	defer rows.Close()

	costumes := []costume.SellerCostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		costume := costume.SellerCostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		costume.Created_at = createdAt.Format("2006-01-02 15:04:05")
		costume.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		costumes = append(costumes, costume)
		hasData = true
	}
	if hasData == false {
		return costumes, errors.New("costume not found")
	}

	return costumes, nil
}

func (repository *CostumeRepository) FindById(ctx context.Context, tx *sql.Tx, id int) (costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at, updated_at FROM costumes where id=$1"
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	costumes := costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	if rows.Next() {
		err := rows.Scan(&costumes.Id, &costumes.User_id, &costumes.Name, &costumes.Description, &costumes.Bahan, &costumes.Ukuran, &costumes.Berat, &costumes.Kategori, &costumes.Price, &costumes.Picture, &costumes.Available, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		costumes.Created_at = createdAt.Format("2006-01-02 15:04:05")
		costumes.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		return costumes, nil
	} else {
		return costumes, errors.New("costume not found")
	}
}

func (repository *CostumeRepository) FindSellerCostumeByCostumeID(ctx context.Context, tx *sql.Tx, userUUID string, costumeID int) (costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at,updated_at FROM costumes WHERE user_id = $1 AND id = $2"
	rows, err := tx.QueryContext(ctx, query, userUUID, costumeID)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	costume := costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	if rows.Next() {
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		costume.Created_at = createdAt.Format("2006-01-02 15:04:05")
		costume.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		return costume, nil
	} else {
		return costume, errors.New("costume not found")
	}
}

func (repository *CostumeRepository) FindPictureById(ctx context.Context, tx *sql.Tx, costumeID int) (*string, error) {
	query := "SELECT costume_picture FROM costumes WHERE id = $1"
	rows, err := tx.QueryContext(ctx, query, costumeID)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	var imagepath *string
	if rows.Next() {
		err = rows.Scan(&imagepath)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return imagepath, nil
	} else {
		return imagepath, errors.New("costume not found")
	}
}

func (repository *CostumeRepository) Delete(ctx context.Context, tx *sql.Tx, id int, uuid string) {
	query := "DELETE FROM costumes WHERE id=$1"
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}
