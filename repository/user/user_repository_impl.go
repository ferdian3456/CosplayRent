package user

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/user"
	"database/sql"
	"errors"
)

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user domain.User) {
	query := "INSERT INTO users (id,name,email,password,role,created_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Role, user.Created_at)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) Login(ctx context.Context, tx *sql.Tx, name string) (domain.User, error) {
	query := "SELECT name,password,role FROM users where name=$1"
	rows, err := tx.QueryContext(ctx, query, name)
	helper.PanicIfError(err)

	defer rows.Close()

	users := domain.User{}
	if rows.Next() {
		err := rows.Scan(&users.Name, &users.Password, &users.Role)
		helper.PanicIfError(err)
		return users, nil
	} else {
		return users, errors.New("wrong name or wrong password")
	}
}

func (repository *UserRepositoryImpl) FindByUUID(ctx context.Context, tx *sql.Tx, uuid string) (user.UserResponse, error) {
	query := "SELECT id,name,email,address,role,profile_picture,created_at FROM users where id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	helper.PanicIfError(err)

	defer rows.Close()

	users := user.UserResponse{}
	if rows.Next() {
		err := rows.Scan(&users.Id, &users.Name, &users.Email, &users.Address, &users.Role, &users.Profile_picture, &users.Created_at)
		helper.PanicIfError(err)
		return users, nil
	} else {
		return users, errors.New("user not found")
	}
}

func (repository *UserRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]user.UserResponse, error) {
	query := "SELECT id,name,email,address,role,profile_picture,created_at FROM users"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	users := []user.UserResponse{}
	for rows.Next() {
		user := user.UserResponse{}
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Address, &user.Role, &user.Profile_picture, &user.Created_at)
		helper.PanicIfError(err)
		users = append(users, user)
		hasData = true
	}
	if hasData == false {
		return users, errors.New("user not found")
	}

	return users, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user user.UserUpdateRequest) {
	query := "UPDATE users SET name=$2,email=$3,password=$4,address=$5,role=$6,profile_picture=$7  WHERE id=$1"
	_, err := tx.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Address, user.Role, user.Profile_picture)
	helper.PanicIfError(err)
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, uuid string) {
	query := "DELETE FROM users WHERE id=$1"
	_, err := tx.ExecContext(ctx, query, uuid)
	helper.PanicIfError(err)
}
