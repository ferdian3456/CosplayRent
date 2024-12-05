package costume

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/costume"
	"database/sql"
	"errors"
	"log"
	"time"
)

type CostumeRepositoryImpl struct{}

func NewCostumeRepository() CostumeRepository {
	return &CostumeRepositoryImpl{}
}

func (repository *CostumeRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, costume domain.Costume) {
	log.Printf("User with uuid: %s enter Costume Repository: Create", costume.User_id)
	query := "INSERT INTO costumes (user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	_, err := tx.ExecContext(ctx, query, costume.User_id, costume.Name, costume.Description, costume.Bahan, costume.Ukuran, costume.Berat, costume.Kategori, costume.Price, costume.Picture, costume.Created_at, costume.Created_at)
	helper.PanicIfError(err)
}

func (repository *CostumeRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at, updated_at FROM costumes where id=$1"
	rows, err := tx.QueryContext(ctx, query, id)
	helper.PanicIfError(err)

	defer rows.Close()

	costumes := costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	if rows.Next() {
		err := rows.Scan(&costumes.Id, &costumes.User_id, &costumes.Name, &costumes.Description, &costumes.Bahan, &costumes.Ukuran, &costumes.Berat, &costumes.Kategori, &costumes.Price, &costumes.Picture, &costumes.Available, &createdAt, &updatedAt)
		helper.PanicIfError(err)
		costumes.Created_at = createdAt.Format("2006-01-02 15:04:05")
		costumes.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		return costumes, nil
	} else {
		return costumes, errors.New("costume not found")
	}
}

func (repository *CostumeRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at, updated_at FROM costumes"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	costumes := []costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		helper.PanicIfError(err)
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

func (repository *CostumeRepositoryImpl) FindByName(ctx context.Context, tx *sql.Tx, name string) ([]costume.CostumeResponse, error) {
	query := "SELECT id,user_id,name,description,price,costume_picture,available,created_at,updated_at FROM costumes WHERE name like $1"
	rows, err := tx.QueryContext(ctx, query, "%"+name+"%")
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	costumes := []costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		helper.PanicIfError(err)
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

func (repository *CostumeRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, costume costume.CostumeUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter Costume Repository: Update", uuid)
	if costume.Picture == nil {
		query := "UPDATE costumes SET name=$2,description=$3,bahan=$4,ukuran=$5,berat=$6,kategori=$7,price=$8,updated_at=$9  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, costume.Id, costume.Name, costume.Description, costume.Bahan, costume.Ukuran, costume.Berat, costume.Kategori, costume.Price, costume.Update_at)
		helper.PanicIfError(err)
	} else {
		query := "UPDATE costumes SET name=$2,description=$3,bahan=$4,ukuran=$5,berat=$6,kategori=$7,price=$8,costume_picture=$9,updated_at=$10  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, costume.Id, costume.Name, costume.Description, costume.Bahan, costume.Ukuran, costume.Berat, costume.Kategori, costume.Price, costume.Picture, costume.Update_at)
		helper.PanicIfError(err)
	}
}

func (repository *CostumeRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id int, uuid string) {
	log.Printf("User with uuid: %s enter Costume Repository: Delete", uuid)
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
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		helper.PanicIfError(err)
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

func (repository *CostumeRepositoryImpl) FindSellerCostumeByCostumeID(ctx context.Context, tx *sql.Tx, userUUID string, costumeID int) (costume.CostumeResponse, error) {
	log.Printf("User with uuid: %s enter Costume Repository: FindSellerCostumeByCostumeID", userUUID)
	//log.Printf("SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at,updated_at FROM costumes WHERE user_id=%s AND id=%d", userUUID, costumeID)
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at,updated_at FROM costumes WHERE user_id = $1 AND id = $2"
	rows, err := tx.QueryContext(ctx, query, userUUID, costumeID)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	costume := costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	if rows.Next() {
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		helper.PanicIfError(err)
		costume.Created_at = createdAt.Format("2006-01-02 15:04:05")
		costume.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		hasData = true
	}
	if hasData == false {
		return costume, errors.New("costume not found")
	}

	return costume, nil
}

func (repository *CostumeRepositoryImpl) FindSellerCostume(ctx context.Context, tx *sql.Tx, uuid string) ([]costume.CostumeResponse, error) {
	log.Printf("User with uuid: %s enter Costume Repository: FindSellerCostume", uuid)
	query := "SELECT id,user_id,name,description,bahan,ukuran,berat,kategori,price,costume_picture,available,created_at, updated_at FROM costumes where user_id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	costumes := []costume.CostumeResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		costume := costume.CostumeResponse{}
		err = rows.Scan(&costume.Id, &costume.User_id, &costume.Name, &costume.Description, &costume.Bahan, &costume.Ukuran, &costume.Berat, &costume.Kategori, &costume.Price, &costume.Picture, &costume.Available, &createdAt, &updatedAt)
		helper.PanicIfError(err)
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

func (repository *CostumeRepositoryImpl) CheckOwnership(ctx context.Context, tx *sql.Tx, userUUID string, costumeid int) error {
	log.Printf("User with uuid: %s enter Costume Repository: CheckOwnership", userUUID)

	query := "SELECT name FROM costumes WHERE id=$1 AND user_id=$2"
	rows, err := tx.QueryContext(ctx, query, costumeid, userUUID)

	helper.PanicIfError(err)
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

func (repository *CostumeRepositoryImpl) GetSellerIdFindByCostumeID(ctx context.Context, tx *sql.Tx, userUUID string, costumeid int) (string, error) {
	log.Printf("User with uuid: %s enter Costume Repository: FindByCostumeId", userUUID)
	query := "SELECT user_id FROM costumes WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, costumeid)

	helper.PanicIfError(err)

	defer row.Close()

	var sellerId string
	if row.Next() {
		err = row.Scan(&sellerId)
		helper.PanicIfError(err)
		return sellerId, nil
	} else {
		return sellerId, errors.New("seller of this costume is not found")
	}
}
