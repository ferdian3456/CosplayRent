package user

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/user"
	"database/sql"
	"errors"
	"log"
	"time"
)

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user domain.User) {
	query := "INSERT INTO users (id,name,email,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Created_at, user.Created_at)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) Login(ctx context.Context, tx *sql.Tx, name string) (domain.User, error) {
	query := "SELECT id,email,password FROM users where email=$1"
	rows, err := tx.QueryContext(ctx, query, name)
	helper.PanicIfError(err)

	defer rows.Close()

	users := domain.User{}
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Email, &users.Password)
		helper.PanicIfError(err)
		return users, nil
	} else {
		return users, errors.New("wrong email or wrong password")
	}
}

func (repository *UserRepositoryImpl) FindByUUID(ctx context.Context, tx *sql.Tx, uuid string) (user.UserResponse, error) {
	log.Printf("User with uuid: %s enter User Repository: FindByUUID", uuid)
	query := "SELECT id,name,email,address,profile_picture,created_at FROM users where id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)

	defer rows.Close()

	users := user.UserResponse{}
	var createdAt time.Time
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Name, &users.Email, &users.Address, &users.Profile_picture, &createdAt)
		helper.PanicIfError(err)
		users.Created_at = createdAt.Format("2006-01-02 15:04:05")
		return users, nil
	} else {
		return users, errors.New("user not found")
	}
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserResponse, error) {
	log.Printf("User with uuid: %s enter User Repository: FindAll", uuid)

	query := "SELECT id,name,email,address,profile_picture,created_at FROM users"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	users := []user.UserResponse{}
	var createdAt time.Time
	for rows.Next() {
		user := user.UserResponse{}
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address, &user.Profile_picture, &createdAt)
		helper.PanicIfError(err)
		user.Created_at = createdAt.Format("2006-01-02 15:04:05")
		//user.Profile_picture = fmt.Sprintf()
		users = append(users, user)
		hasData = true
	}
	if hasData == false {
		return users, errors.New("user not found")
	}

	return users, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user user.UserUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter User Repository: Update", uuid)

	if user.Profile_picture == nil {
		query := "UPDATE users SET name=$2,email=$3,address=$4,updated_at=$5  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Address, user.Update_at)
		helper.PanicIfError(err)
	} else {
		query := "UPDATE users SET name=$2,email=$3,address=$4,profile_picture=$5,updated_at=$6  WHERE id=$1"
		_, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Address, user.Profile_picture, user.Update_at)
		helper.PanicIfError(err)
	}
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, uuid string) {
	log.Printf("User with uuid: %s enter User Repository: Delete", uuid)
	query := "DELETE FROM users WHERE id=$1"
	_, err := tx.ExecContext(ctx, query, uuid)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (user.UserResponse, error) {
	query := "SELECT id,name,email,address,profile_picture,created_at FROM users where email=$1"
	rows, err := tx.QueryContext(ctx, query, email)
	helper.PanicIfError(err)

	defer rows.Close()

	users := user.UserResponse{}
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Name, &users.Email, &users.Address, &users.Profile_picture, &users.Created_at)
		helper.PanicIfError(err)
		return users, nil
	} else {
		return users, errors.New("user not found")
	}
}

func (repository *UserRepositoryImpl) AddOrUpdateIdentityCard(ctx context.Context, tx *sql.Tx, uuid string, IdentityCardImage string) {
	log.Printf("User with uuid: %s enter User Repository: AddOrUpdateIdentityCard", uuid)

	query := "UPDATE users SET identitycard_picture = $1 WHERE id = $2"
	_, err := tx.ExecContext(ctx, query, IdentityCardImage, uuid)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) GetIdentityCard(ctx context.Context, tx *sql.Tx, uuid string) (string, error) {
	log.Printf("User with uuid: %s enter User Repository: GetIdentityCard", uuid)

	query := "SELECT identitycard_picture FROM users WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)

	defer row.Close()

	var IdentityCardImage *string
	if row.Next() {
		err := row.Scan(&IdentityCardImage)
		helper.PanicIfError(err)
		if IdentityCardImage != nil {
			return *IdentityCardImage, nil
		} else {
			return "", errors.New("identity card is empty")
		}
	} else {
		return "", errors.New("identity card not found")
	}
}

func (repository *UserRepositoryImpl) GetEMoneyAmount(ctx context.Context, tx *sql.Tx, uuid string) (float64, error) {
	log.Printf("User with uuid: %s enter User Repository: GetEMoneyAmount", uuid)

	query := "SELECT emoney_amount FROM users WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)

	defer row.Close()

	var eMoneyAmount float64

	if row.Next() {
		err = row.Scan(&eMoneyAmount)
		helper.PanicIfError(err)
		return eMoneyAmount, nil
	} else {
		return 0, errors.New("Emoney Amount is Not Found")
	}
}

func (repository *UserRepositoryImpl) TopUp(ctx context.Context, tx *sql.Tx, emoney user.TopUpEmoney, uuid string) {
	log.Printf("User with uuid: %s enter User Repository: Topup", uuid)

	query := "UPDATE users SET emoney_amount = $1 WHERE id = $2"
	_, err := tx.ExecContext(ctx, query, emoney.Emoney_amont, uuid)
	helper.PanicIfError(err)
}
