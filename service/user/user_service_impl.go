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
		"name":    userDomain.Name,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func (service *UserServiceImpl) Login(ctx context.Context, request user.UserLoginRequest) string {
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
	helper.PanicIfError(err)

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	secretKey := os.Getenv("SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func (service *UserServiceImpl) FindByUUID(ctx context.Context, uuid string) user.UserResponse {
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

	return user
}

func (service *UserServiceImpl) FindAll(ctx context.Context) []user.UserResponse {
	var err error
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToUserResponses(user)
}

func (service *UserServiceImpl) Update(ctx context.Context, userRequest user.UserUpdateRequest) {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	result, err1 := service.UserRepository.FindByUUID(ctx, tx, userRequest.Id)
	if err1 != nil {
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

	service.UserRepository.Update(ctx, tx, updateRequest)
}

func (service *UserServiceImpl) Delete(ctx context.Context, uuid string) {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	service.UserRepository.Delete(ctx, tx, uuid)
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

	var email string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if val, exists := claims["email"]; exists {
			if strVal, ok := val.(string); ok {
				email = strVal
			} else {
				return user.UserResponse{}, fmt.Errorf("name claim is not a string")
			}
		} else {
			return user.UserResponse{}, fmt.Errorf("name claim does not exist")
		}
	}

	//log.Println(name)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)
	userDomain, err := service.UserRepository.FindByEmail(ctx, tx, email)
	helper.PanicIfError(err)

	return userDomain, nil
}
