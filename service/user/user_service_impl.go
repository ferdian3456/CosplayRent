package user

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/user"
	users "cosplayrent/repository/user"
	"database/sql"
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

func (service *UserServiceImpl) Create(ctx context.Context, request user.UserCreateRequest) {
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
		Role:       request.Role,
		Created_at: &now,
	}

	service.UserRepository.Create(ctx, tx, userDomain)
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
		Name:     request.Name,
		Password: request.Password,
	}

	user := domain.User{}
	user, err = service.UserRepository.Login(ctx, tx, userDomain.Name)
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
		"name":    userDomain.Name,
		"role":    userDomain.Role,
		"expired": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
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

	updateRequest := user.UserUpdateRequest{
		Id:              result.Id,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Address:         userRequest.Address,
		Password:        userRequest.Password,
		Role:            userRequest.Role,
		Profile_picture: userRequest.Profile_picture,
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
