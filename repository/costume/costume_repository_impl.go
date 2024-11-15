package costume

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/costume"
	"database/sql"
	"errors"
)

type CostumeRepositoryImpl struct{}

func NewCostumeRepository() CostumeRepository {
	return &CostumeRepositoryImpl{}
}

func (repository *CostumeRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, costume domain.Costume) {
	query := "INSERT INTO costumes (user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	_, err := tx.ExecContext(ctx, query, costume.User_id, costume.Name, costume.Description, costume.Bahan, costume.Ukuran, costume.Berat, costume.Kategori, costume.Price, costume.Picture, costume.Available, costume.Created_at)
	helper.PanicIfError(err)
}

func (repository *CostumeRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at FROM costumes where id=$1"
	rows, err := tx.QueryContext(ctx, query, id)
	helper.PanicIfError(err)

	defer rows.Close()

	costumes := costume.CostumeResponse{}
	if rows.Next() {
		err := rows.Scan(&costumes.Id, &costumes.User_id, &costumes.Name, &costumes.Description, &costumes.Bahan, &costumes.Ukuran, &costumes.Berat, &costumes.Kategori, &costumes.Price, &costumes.Picture, &costumes.Available, &costumes.Created_at)
		helper.PanicIfError(err)
		return costumes, nil
	} else {
		return costumes, errors.New("costume not found")
	}
}

func (repository *CostumeRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at FROM costumes"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	costumes := []costume.CostumeResponse{}
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &costume.Created_at)
		helper.PanicIfError(err)
		costumes = append(costumes, costume)
		hasData = true
	}
	if hasData == false {
		return costumes, errors.New("costume not found")
	}

	return costumes, nil
}

func (repository *CostumeRepositoryImpl) FindByName(ctx context.Context, tx *sql.Tx, name string) ([]costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,price,costume_picture,available,created_at FROM costumes WHERE name like $1"
	rows, err := tx.QueryContext(ctx, query, "%"+name+"%")
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	costumes := []costume.CostumeResponse{}
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Price, &costume.Picture, &costume.Available, &costume.Created_at)
		helper.PanicIfError(err)
		costumes = append(costumes, costume)
		hasData = true
	}
	if hasData == false {
		return costumes, errors.New("costume not found")
	}

	return costumes, nil
}

func (repository *CostumeRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, costume costume.CostumeUpdateRequest) {
	if costume.Picture == nil {
		query := "UPDATE costumes SET name=$2,description=$3,bahan=$4,ukuran=$5,berat=$6,kategori=$7,price=$8,available=$9,updated_at=$10  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, costume.Id, costume.Name, costume.Description, costume.Bahan, costume.Ukuran, costume.Berat, costume.Kategori, costume.Price, costume.Available, costume.Update_at)
		helper.PanicIfError(err)
	} else {
		query := "UPDATE costumes SET name=$2,description=$3,bahan=$4,ukuran=$5,berat=$6,kategori=$7,price=$8,costume_picture=$9,available=$10,updated_at=$11  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, costume.Id, costume.Name, costume.Description, costume.Bahan, costume.Ukuran, costume.Berat, costume.Kategori, costume.Price, costume.Picture, costume.Available, costume.Update_at)
		helper.PanicIfError(err)
	}
}

func (repository *CostumeRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int) {
	query := "DELETE FROM costumes WHERE id=$1"
	_, err := tx.ExecContext(ctx, query, id)
	helper.PanicIfError(err)
}

func (repository *CostumeRepositoryImpl) FindByUserUUID(ctx context.Context, tx *sql.Tx, userUUID string) ([]costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at, updated_at FROM costumes WHERE user_id=$1"
	rows, err := tx.QueryContext(ctx, query, userUUID)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	costumes := []costume.CostumeResponse{}
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &costume.Created_at, &costume.Updated_at)
		helper.PanicIfError(err)
		costumes = append(costumes, costume)
		hasData = true
	}
	if hasData == false {
		return costumes, errors.New("costume not found")
	}

	return costumes, nil
}

func (repository *CostumeRepositoryImpl) FindSellerCostumeByCostumeID(ctx context.Context, tx *sql.Tx, userUUID string, costumeID int) (costume.CostumeResponse, error) {
	query := "SELECT id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at,updated_at WHERE user_id=$1 AND id=$2"
	rows, err := tx.QueryContext(ctx, query, userUUID, costumeID)
	helper.PanicIfError(err)
	hasData := false

	costumes := costume.CostumeResponse{}
	if rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &costume.Created_at, &costume.Updated_at)
		helper.PanicIfError(err)
		hasData = true
	}
	if hasData == false {
		return costumes, errors.New("costume not found")
	}

	return costumes, nil
}
