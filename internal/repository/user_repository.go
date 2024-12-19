package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/user"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type UserRepository struct {
	Log *zerolog.Logger
}

func NewUserRepository(zerolog *zerolog.Logger) *UserRepository {
	return &UserRepository{
		Log: zerolog,
	}
}

func (repository *UserRepository) Create(ctx context.Context, tx *sql.Tx, user domain.User) {
	query := "INSERT INTO users (id,name,email,password,emoney_updated_at,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Created_at, user.Created_at, user.Created_at)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(respErr).Msg(err.Error())
	}
}

func (repository *UserRepository) CheckCredentialUnique(ctx context.Context, tx *sql.Tx, user domain.User) error {
	query := "SELECT name, email from users WHERE name=$1 OR email=$2 LIMIT 1"
	row, err := tx.QueryContext(ctx, query, user.Name, user.Email)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer row.Close()

	if row.Next() {
		return errors.New("name or email are already exists")
	} else {
		return nil
	}
}

func (repository *UserRepository) Login(ctx context.Context, tx *sql.Tx, name string) (domain.User, error) {
	query := "SELECT id,email,password FROM users where email=$1"
	rows, err := tx.QueryContext(ctx, query, name)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer rows.Close()

	users := domain.User{}
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Email, &users.Password)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(respErr).Msg(err.Error())
		}
		return users, nil
	} else {
		return users, errors.New("wrong email")
	}
}

func (repository *UserRepository) FindByUUID(ctx context.Context, tx *sql.Tx, uuid string) (user.UserResponse, error) {
	query := "SELECT id,name,email,address,profile_picture,originprovince_name,originprovince_id,origincity_name,origincity_id,created_at,updated_at FROM users where id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer rows.Close()

	users := user.UserResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Name, &users.Email, &users.Address, &users.Profile_picture, &users.Origin_province_name, &users.Origin_province_id, &users.Origin_city_name, &users.Origin_city_id, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(respErr).Msg(err.Error())
		}
		users.Created_at = createdAt.Format("2006-01-02 15:04:05")
		users.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		return users, nil
	} else {
		return users, errors.New("user not found")
	}
}

func (repository *UserRepository) CheckUserExistance(ctx context.Context, tx *sql.Tx, uuid string) error {
	query := "SELECT id FROM users where id=$1"
	row, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer row.Close()

	if row.Next() {
		return nil
	} else {
		return errors.New("user not found")
	}
}

func (repository *UserRepository) FindAll(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserResponse, error) {
	query := "SELECT id,name,email,address,profile_picture,originprovince_name,originprovince_id,origincity_name,origincity_id,created_at,updated_at FROM users"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(respErr).Msg(err.Error())
	}

	hasData := false

	defer rows.Close()

	users := []user.UserResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	for rows.Next() {
		user := user.UserResponse{}
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address, &user.Profile_picture, &user.Origin_province_name, &user.Origin_province_id, &user.Origin_city_name, &user.Origin_city_id, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(respErr).Msg(err.Error())
		}
		user.Created_at = createdAt.Format("2006-01-02 15:04:05")
		user.Updated_at = updatedAt.Format("2006-01-02 15:04:05")
		users = append(users, user)
		hasData = true
	}
	if hasData == false {
		return users, errors.New("user not found")
	}

	return users, nil
}

func (repository *UserRepository) Update(ctx context.Context, tx *sql.Tx, user domain.User) {
	setClauses := []string{}
	args := []interface{}{user.Id}

	if user.Name != "" {
		setClauses = append(setClauses, "name = $2")
		args = append(args, user.Name)
	}

	if user.Email != "" {
		setClauses = append(setClauses, "email = $3")
		args = append(args, user.Email)
	}

	if user.Profile_picture != "" {
		setClauses = append(setClauses, "profile_picture = $4")
		args = append(args, user.Profile_picture)
	}

	if user.Address != "" {
		setClauses = append(setClauses, "address = $5")
		args = append(args, user.Address)
	}

	if user.Origin_province_name != "" {
		setClauses = append(setClauses, "originprovince_name = $6")
		args = append(args, user.Origin_province_name)
	}

	if user.Origin_province_id != 0 {
		setClauses = append(setClauses, "originprovince_id = $7")
		args = append(args, user.Origin_province_id)
	}

	if user.Origin_city_name != "" {
		setClauses = append(setClauses, "origincity_name = $8")
		args = append(args, user.Origin_city_name)
	}

	if user.Origin_city_id != 0 {
		setClauses = append(setClauses, "origincity_id = $9")
		args = append(args, user.Origin_city_id)
	}

	setClauses = append(setClauses, "updated_at = $10")
	args = append(args, user.Updated_at)

	query := fmt.Sprintf(`
        UPDATE users 
        SET %s
        WHERE id = $1;
    `, strings.Join(setClauses, ", "))

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(respErr).Msg(err.Error())
	}
}

//
//func (repository *UserRepository) Delete(ctx context.Context, tx *sql.Tx, uuid string) {
//	log.Printf("User with uuid: %s enter User Repository: Delete", uuid)
//	query := "DELETE FROM users WHERE id=$1"
//	_, err := tx.ExecContext(ctx, query, uuid)
//	helper.PanicIfError(err)
//}
//
//func (repository *UserRepository) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (user.UserResponse, error) {
//	query := "SELECT id,name,email,address,profile_picture,created_at FROM users where email=$1"
//	rows, err := tx.QueryContext(ctx, query, email)
//	helper.PanicIfError(err)
//
//	defer rows.Close()
//
//	users := user.UserResponse{}
//	if rows.Next() {
//		err := rows.Scan(&users.Id, &users.Name, &users.Email, &users.Address, &users.Profile_picture, &users.Created_at)
//		helper.PanicIfError(err)
//		return users, nil
//	} else {
//		return users, errors.New("user not found")
//	}
//}
//
//func (repository *UserRepository) AddOrUpdateIdentityCard(ctx context.Context, tx *sql.Tx, uuid string, IdentityCardImage string) {
//	log.Printf("User with uuid: %s enter User Repository: AddOrUpdateIdentityCard", uuid)
//
//	if IdentityCardImage != "" {
//		query := "UPDATE users SET identitycard_picture = $1 WHERE id = $2"
//		_, err := tx.ExecContext(ctx, query, IdentityCardImage, uuid)
//		helper.PanicIfError(err)
//	}
//}
//
//func (repository *UserRepository) GetIdentityCard(ctx context.Context, tx *sql.Tx, uuid string) (string, error) {
//	log.Printf("User with uuid: %s enter User Repository: GetIdentityCard", uuid)
//
//	query := "SELECT identitycard_picture FROM users WHERE id=$1"
//	row, err := tx.QueryContext(ctx, query, uuid)
//	helper.PanicIfError(err)
//
//	defer row.Close()
//
//	var IdentityCardImage *string
//	if row.Next() {
//		err := row.Scan(&IdentityCardImage)
//		helper.PanicIfError(err)
//		if IdentityCardImage != nil {
//			return *IdentityCardImage, nil
//		} else {
//			return "", errors.New("identity card is empty")
//		}
//	} else {
//		return "", errors.New("identity card not found")
//	}
//}
//
//func (repository *UserRepository) GetEMoneyAmount(ctx context.Context, tx *sql.Tx, uuid string) (user.UserEmoneyResponse, error) {
//	log.Printf("User with uuid: %s enter User Repository: GetEMoneyAmount", uuid)
//
//	query := "SELECT emoney_amount,emoney_updated_at FROM users WHERE id=$1"
//	row, err := tx.QueryContext(ctx, query, uuid)
//	helper.PanicIfError(err)
//
//	defer row.Close()
//
//	userEmoney := user.UserEmoneyResponse{}
//	var updatedAt time.Time
//	if row.Next() {
//		err = row.Scan(&userEmoney.Emoney_amont, &updatedAt)
//		helper.PanicIfError(err)
//		userEmoney.Emoney_updated_at = updatedAt.Format("2006-01-02 15:04:05")
//		return userEmoney, nil
//	} else {
//		return user.UserEmoneyResponse{}, errors.New("emoney amount is not found")
//	}
//}
//
//func (repository *UserRepository) TopUp(ctx context.Context, tx *sql.Tx, emoney float64, uuid string, timeNow *time.Time) {
//	log.Printf("User with uuid: %s enter User Repository: Topup", uuid)
//
//	query := "UPDATE users SET emoney_amount = emoney_amount + $1, emoney_updated_at=$2 WHERE id = $3"
//	_, err := tx.ExecContext(ctx, query, emoney, timeNow, uuid)
//	helper.PanicIfError(err)
//}
//
//func (repository *UserRepository) AfterBuy(ctx context.Context, tx *sql.Tx, orderamount float64, buyeruuid string, selleruuid string, timeNow *time.Time) {
//	log.Printf("Buy with uuid: %s and Seller with uuid: %s enter User Repository: AfterBuy", buyeruuid, selleruuid)
//
//	//log.Println("TimeNow:", timeNow)
//	// substract buyer money
//	query := "UPDATE users SET emoney_amount = emoney_amount - $1,emoney_updated_at=$2 WHERE id = $3"
//	_, err := tx.ExecContext(ctx, query, orderamount, timeNow, buyeruuid)
//
//	helper.PanicIfError(err)
//
//	// add seller money
//
//	query = "UPDATE users SET emoney_amount = emoney_amount + $1, emoney_updated_at=$2 WHERE id = $3"
//	_, err = tx.ExecContext(ctx, query, orderamount, timeNow, selleruuid)
//
//	helper.PanicIfError(err)
//}
//
//func (repository *UserRepository) CheckUserStatus(ctx context.Context, tx *sql.Tx, userid string) (user.CheckUserStatusResponse, error) {
//	log.Printf("User with uuid: %s enter User Repository: CheckUserStatus", userid)
//
//	query := "SELECT id,name,identitycard_picture,address,origincity_name FROM users WHERE id=$1"
//	row, err := tx.QueryContext(ctx, query, userid)
//	helper.PanicIfError(err)
//
//	defer row.Close()
//
//	checkuserStatus := user.CheckUserStatusResponse{}
//	var IdentityCardImage *string
//	var originCityName *string
//	var address *string
//	if row.Next() {
//		err := row.Scan(&checkuserStatus.User_id, &checkuserStatus.Name, &IdentityCardImage, &address, &originCityName)
//		helper.PanicIfError(err)
//		if IdentityCardImage != nil && originCityName != nil && address != nil {
//			return checkuserStatus, nil
//		} else {
//			return checkuserStatus, errors.New("need to fulfill identity card and address detail (address,province, and city)")
//		}
//	} else {
//		return checkuserStatus, errors.New("user not found")
//	}
//}
