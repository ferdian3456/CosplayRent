package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/user"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
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
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *UserRepository) CheckCredentialUnique(ctx context.Context, tx *sql.Tx, user domain.User) error {
	query := "SELECT name, email from users WHERE name=$1 OR email=$2 LIMIT 1"
	row, err := tx.QueryContext(ctx, query, user.Name, user.Email)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
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
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	users := domain.User{}
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Email, &users.Password)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
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
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	users := user.UserResponse{}
	var createdAt time.Time
	var updatedAt time.Time
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Name, &users.Email, &users.Address, &users.Profile_picture, &users.Origin_province_name, &users.Origin_province_id, &users.Origin_city_name, &users.Origin_city_id, &createdAt, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
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
		repository.Log.Panic().Err(err).Msg(respErr.Error())
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
		repository.Log.Panic().Err(err).Msg(respErr.Error())
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
			repository.Log.Panic().Err(err).Msg(respErr.Error())
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
	query := "UPDATE users SET "
	args := []interface{}{}
	argCounter := 1

	if user.Name != "" {
		query += fmt.Sprintf("name = $%d, ", argCounter)
		args = append(args, user.Name)
		argCounter++
	}
	if user.Email != "" {
		query += fmt.Sprintf("email = $%d, ", argCounter)
		args = append(args, user.Email)
		argCounter++
	}
	if user.Profile_picture != "" {
		query += fmt.Sprintf("profile_picture = $%d, ", argCounter)
		args = append(args, user.Profile_picture)
		argCounter++
	}
	if user.Address != "" {
		query += fmt.Sprintf("address = $%d, ", argCounter)
		args = append(args, user.Address)
		argCounter++
	}
	if user.Origin_province_name != "" {
		query += fmt.Sprintf("originprovince_name = $%d, ", argCounter)
		args = append(args, user.Origin_province_name)
		argCounter++
	}
	if user.Origin_province_id != 0 {
		query += fmt.Sprintf("originprovince_id = $%d, ", argCounter)
		args = append(args, user.Origin_province_id)
		argCounter++
	}
	if user.Origin_city_name != "" {
		query += fmt.Sprintf("origincity_name = $%d, ", argCounter)
		args = append(args, user.Origin_city_name)
		argCounter++
	}
	if user.Origin_city_id != 0 {
		query += fmt.Sprintf("origincity_id = $%d, ", argCounter)
		args = append(args, user.Origin_city_id)
		argCounter++
	}

	query += fmt.Sprintf("updated_at = $%d ", argCounter)
	args = append(args, user.Updated_at)
	argCounter++

	query += fmt.Sprintf("WHERE id = $%d", argCounter)
	args = append(args, user.Id)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *UserRepository) AddOrUpdateIdentityCard(ctx context.Context, tx *sql.Tx, user domain.User) {
	if user.Identity_card_picture != "" {
		query := "UPDATE users SET identitycard_picture = $1 WHERE id = $2"
		_, err := tx.ExecContext(ctx, query, user.Identity_card_picture, user.Id)
		if err != nil {
			respErr := errors.New("failed to query into database")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}
}

func (repository *UserRepository) GetIdentityCard(ctx context.Context, tx *sql.Tx, uuid string) (string, error) {
	query := "SELECT identitycard_picture FROM users WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var IdentityCardImage *string
	if row.Next() {
		err := row.Scan(&IdentityCardImage)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		if IdentityCardImage != nil {
			return *IdentityCardImage, nil
		} else {
			return "", errors.New("identity card is empty")
		}
	} else {
		return "", errors.New("identity card not found")
	}
}

func (repository *UserRepository) GetEMoneyAmount(ctx context.Context, tx *sql.Tx, uuid string) user.UserEmoneyResponse {
	query := "SELECT emoney_amount,emoney_updated_at FROM users WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	user := user.UserEmoneyResponse{}
	var updatedAt time.Time
	if row.Next() {
		err = row.Scan(&user.Emoney_amont, &updatedAt)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		user.Emoney_updated_at = updatedAt.Format("2006-01-02 15:04:05")
		return user
	} else {
		user.Emoney_amont = 0
		return user
	}
}

func (repository *UserRepository) FindAllMoneyChanges(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserEMoneyTransactionHistory, error) {
	query := `
    SELECT total AS amount, updated_at, 'Order (Buyer)' AS source
    FROM orders
    WHERE user_id = $1 AND status_payment = true

    UNION ALL

    SELECT total AS amount, updated_at, 'Order (Seller)' AS source
    FROM orders
    WHERE seller_id = $1 AND status_payment = true

    UNION ALL

    SELECT topup_amount AS amount, updated_at, 'Top Up' AS source
    FROM topup_orders
    WHERE user_id = $1 AND status_payment = true;
	`
	rows, err := tx.QueryContext(ctx, query, uuid)
	hasData := false

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	users := []user.UserEMoneyTransactionHistory{}
	var updatedAt time.Time
	var transactionType string

	for rows.Next() {
		user := user.UserEMoneyTransactionHistory{}
		err = rows.Scan(&user.Transaction_amount, &updatedAt, &transactionType)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		user.Transaction_date = updatedAt.Format("2006-01-02 15:04:05")
		user.Transaction_type = transactionType
		hasData = true
		users = append(users, user)
	}
	if hasData == false {
		return users, errors.New("transaction history not found")
	}

	return users, nil
}

func (repository *UserRepository) FindNameAndEmailById(ctx context.Context, tx *sql.Tx, uuid string) (domain.User, error) {
	query := "SELECT name,email FROM users WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	user := domain.User{}
	if row.Next() {
		err := row.Scan(&user.Name, &user.Email)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}

func (repository *UserRepository) TopUp(ctx context.Context, tx *sql.Tx, midtrans domain.Midtrans) {
	query := "UPDATE users SET emoney_amount = emoney_amount + $1, emoney_updated_at=$2 WHERE id = $3"
	_, err := tx.ExecContext(ctx, query, midtrans.Order_amount, midtrans.Updated_at, midtrans.TopUpUser_id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *UserRepository) AfterBuy(ctx context.Context, tx *sql.Tx, midtrans domain.Midtrans) {
	query := "UPDATE users SET emoney_amount = emoney_amount - $1,emoney_updated_at=$2 WHERE id = $3"
	_, err := tx.ExecContext(ctx, query, midtrans.Order_amount, midtrans.Updated_at, midtrans.OrderBuyer_id)

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	query = "UPDATE users SET emoney_amount = emoney_amount + $1, emoney_updated_at=$2 WHERE id = $3"
	_, err = tx.ExecContext(ctx, query, midtrans.Order_amount, midtrans.Updated_at, midtrans.OrderSeller_id)

	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *UserRepository) CheckUserStatus(ctx context.Context, tx *sql.Tx, userid string) (user.CheckUserStatusResponse, error) {
	query := "SELECT id,name,identitycard_picture,address,origincity_name FROM users WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, userid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	checkuserStatus := user.CheckUserStatusResponse{}
	var IdentityCardImage *string
	var originCityName *string
	var address *string
	if row.Next() {
		err := row.Scan(&checkuserStatus.User_id, &checkuserStatus.Name, &IdentityCardImage, &address, &originCityName)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		if IdentityCardImage != nil && originCityName != nil && address != nil {
			return checkuserStatus, nil
		} else {
			return checkuserStatus, errors.New("need to fulfill identity card and address detail (address,province, and city)")
		}
	} else {
		return checkuserStatus, errors.New("user not found")
	}
}

func (repository *UserRepository) FindAddressByUserId(ctx context.Context, tx *sql.Tx, userid string) (domain.User, error) {
	query := "SELECT name,originprovince_name,originprovince_id,origincity_name,origincity_id FROM users WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, userid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	user := domain.User{}

	if row.Next() {
		err := row.Scan(&user.Name, &user.Origin_province_name, &user.Origin_province_id, &user.Origin_city_name, &user.Origin_city_id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return user, nil
	} else {
		return user, errors.New("user not found")
	}
}
