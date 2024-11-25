package user

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/user"
	users "cosplayrent/repository/user"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type UserServiceImpl struct {
	UserRepository users.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository users.UserRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) Create(ctx context.Context, request user.UserCreateRequest) string {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	helper.PanicIfError(err)

	userDomain := domain.User{
		Id:         uuid.New(),
		Name:       request.Name,
		Email:      request.Email,
		Password:   string(hashedPassword),
		Created_at: &now,
	}

	service.UserRepository.Create(ctx, tx, userDomain)

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	secretKey := os.Getenv("SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      userDomain.Id,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func (service *UserServiceImpl) Login(ctx context.Context, request user.UserLoginRequest) (string, error) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	userDomain := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user := domain.User{}
	user, err = service.UserRepository.Login(ctx, tx, userDomain.Email)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDomain.Password))
	if err != nil {
		return "", err
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	secretKey := os.Getenv("SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.Id,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		panic(err)
	}

	return tokenString, nil
}

func (service *UserServiceImpl) FindByUUID(ctx context.Context, uuid string) user.UserResponse {
	log.Printf("User with uuid: %s enter User Service: FindByUUID", uuid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	user := user.UserResponse{}
	user, err = service.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	if user.Profile_picture != nil {
		value := imageEnv + *user.Profile_picture
		user.Profile_picture = &value
	}

	return user
}

func (service *UserServiceImpl) FindAll(ctx context.Context, uuid string) []user.UserResponse {
	log.Printf("User with uuid: %s enter User Controller: FindAll", uuid)

	var err error
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindAll(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range user {
		if user[i].Profile_picture != nil {
			value := imageEnv + *user[i].Profile_picture
			user[i].Profile_picture = &value
		}
	}

	return helper.ToUserResponses(user)
}

func (service *UserServiceImpl) Update(ctx context.Context, userRequest user.UserUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter User Service: Update", uuid)

	err := service.Validate.Struct(userRequest)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	result, err := service.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	now := time.Now()

	updateRequest := user.UserUpdateRequest{
		Id:              result.Id,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Address:         userRequest.Address,
		Profile_picture: userRequest.Profile_picture,
		Update_at:       &now,
	}

	service.UserRepository.Update(ctx, tx, updateRequest, uuid)
}

func (service *UserServiceImpl) Delete(ctx context.Context, uuid string) {
	log.Printf("User with uuid: %s enter User Service: Delete", uuid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	userResult, err := service.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.UserRepository.Delete(ctx, tx, uuid)

	finalProfilePicturePath := ".." + *userResult.Profile_picture

	err = os.Remove(finalProfilePicturePath)
	helper.PanicIfError(err)
}

func (service *UserServiceImpl) VerifyAndRetrieve(ctx context.Context, tokenString string) (user.UserResponse, error) {
	secretKey := os.Getenv("SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKeyByte, nil
	})

	if err != nil || !token.Valid {
		return user.UserResponse{}, errors.New("token is not valid")
	}

	var id string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if val, exists := claims["id"]; exists {
			if strVal, ok := val.(string); ok {
				id = strVal
			} else {
				return user.UserResponse{}, fmt.Errorf("id claim is not a string")
			}
		} else {
			return user.UserResponse{}, fmt.Errorf("id claim does not exist")
		}
	}

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)
	userDomain, err := service.UserRepository.FindByUUID(ctx, tx, id)
	helper.PanicIfError(err)

	return userDomain, nil
}
