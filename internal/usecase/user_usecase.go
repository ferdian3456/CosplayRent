package usecase

import (
	"context"
	"cosplayrent/internal/exception"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	googleuuid "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepository *repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
	Log            *zerolog.Logger
}

func NewUserUsecase(userRepository *repository.UserRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger) *UserUsecase {
	return &UserUsecase{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
		Log:            zerolog,
	}
}

func (Usecase *UserUsecase) Create(ctx context.Context, request user.UserCreateRequest) string {
	err := Usecase.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	helper.PanicIfError(err)

	userDomain := domain.User{
		Id:         googleuuid.New(),
		Name:       request.Name,
		Email:      request.Email,
		Password:   string(hashedPassword),
		Created_at: &now,
	}

	Usecase.UserRepository.Create(ctx, tx, userDomain)

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

func (Usecase *UserUsecase) Login(ctx context.Context, request user.UserLoginRequest) (string, error) {
	err := Usecase.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	userDomain := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user := domain.User{}
	user, err = Usecase.UserRepository.Login(ctx, tx, userDomain.Email)
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

func (Usecase *UserUsecase) FindByUUID(ctx context.Context, uuid string) user.UserResponse {
	log.Printf("User with uuid: %s enter User Usecase: FindByUUID", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	user := user.UserResponse{}
	user, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
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

func (Usecase *UserUsecase) FindAll(ctx context.Context, uuid string) []user.UserResponse {
	log.Printf("User with uuid: %s enter User Controller: FindAll", uuid)

	var err error
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	user, err := Usecase.UserRepository.FindAll(ctx, tx, uuid)
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

	return user
}

func (Usecase *UserUsecase) Update(ctx context.Context, userRequest user.UserUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter User Usecase: Update", uuid)

	err := Usecase.Validate.Struct(userRequest)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	result, err := Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	now := time.Now()

	updateRequest := user.UserUpdateRequest{
		Id:                   result.Id,
		Name:                 userRequest.Name,
		Email:                userRequest.Email,
		Address:              userRequest.Address,
		Profile_picture:      userRequest.Profile_picture,
		Origin_province_name: userRequest.Origin_province_name,
		Origin_province_id:   userRequest.Origin_province_id,
		Origin_city_name:     userRequest.Origin_city_name,
		Origin_city_id:       userRequest.Origin_city_id,
		Update_at:            &now,
	}

	Usecase.UserRepository.Update(ctx, tx, updateRequest, uuid)
}

func (Usecase *UserUsecase) Delete(ctx context.Context, uuid string) {
	log.Printf("User with uuid: %s enter User Usecase: Delete", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	Usecase.UserRepository.Delete(ctx, tx, uuid)

	finalProfilePicturePath := ".." + *userResult.Profile_picture

	err = os.Remove(finalProfilePicturePath)
	helper.PanicIfError(err)
}

func (Usecase *UserUsecase) VerifyAndRetrieve(ctx context.Context, tokenString string) (user.UserResponse, error) {
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

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)
	userDomain, err := Usecase.UserRepository.FindByUUID(ctx, tx, id)
	helper.PanicIfError(err)

	return userDomain, nil
}

func (Usecase *UserUsecase) AddIdentityCard(ctx context.Context, uuid string, IdentityCardImage string) {
	log.Printf("User with uuid: %s enter User Usecase: AddIdentityCard", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	Usecase.UserRepository.AddOrUpdateIdentityCard(ctx, tx, uuid, IdentityCardImage)
}

func (Usecase *UserUsecase) GetIdentityCard(ctx context.Context, uuid string) (identityCardImage string) {
	log.Printf("User with uuid: %s enter User Usecase: GetIdentityCard", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	identityCardResult, err := Usecase.UserRepository.GetIdentityCard(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	identityCardResult = imageEnv + identityCardResult

	return identityCardResult
}

func (Usecase *UserUsecase) UpdateIdentityCard(ctx context.Context, uuid string, IdentityCardImage string) {
	log.Printf("User with uuid: %s enter User Usecase: UpdateIdentityCard", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	Usecase.UserRepository.AddOrUpdateIdentityCard(ctx, tx, uuid, IdentityCardImage)
}

func (Usecase *UserUsecase) GetEMoneyAmount(ctx context.Context, uuid string) (userEmoneyResult user.UserEmoneyResponse) {
	log.Printf("User with uuid: %s enter User Usecase: GetEMoneyAmount", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userEmoneyResult, err = Usecase.UserRepository.GetEMoneyAmount(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return userEmoneyResult
}

func (Usecase *UserUsecase) GetEMoneyTransactionHistory(ctx context.Context, uuid string) []user.UserEMoneyTransactionHistory {
	log.Printf("User with uuid: %s enter User Usecase: GetEMoneyTransactionHistory", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	orderHistoryResult, err1 := Usecase.OrderRepository.FindOrderHistoryByUserId(ctx, tx, uuid)

	orderHistorySellerResult, err := Usecase.OrderRepository.FindOrderHistoryBySellerId(ctx, tx, uuid)

	topupOrderHistoryResult, err2 := Usecase.TopUpOrderRepository.FindTopUpOrderHistoryByUserId(ctx, tx, uuid)

	err3 := errors.New("order and topup order is not found")

	if err != nil && err1 != nil && err2 != nil {
		panic(exception.NewNotFoundError(err3.Error()))
	}

	EMoneyOrderHistory := []user.UserEMoneyTransactionHistory{}
	EMoneyOrderSellerHistory := []user.UserEMoneyTransactionHistory{}
	EMoneyTopUpOrderHistory := []user.UserEMoneyTransactionHistory{}

	for _, order := range orderHistoryResult {
		EMoneyOrderHistory = append(EMoneyOrderHistory, user.UserEMoneyTransactionHistory{
			Transaction_amount: order.Emoney_amont,
			Transaction_type:   "Order (Buyer)",
			Transaction_date:   order.Emoney_updated_at,
		})
	}

	for _, order := range orderHistorySellerResult {
		EMoneyOrderSellerHistory = append(EMoneyOrderSellerHistory, user.UserEMoneyTransactionHistory{
			Transaction_amount: order.Emoney_amont,
			Transaction_type:   "Order (Seller)",
			Transaction_date:   order.Emoney_updated_at,
		})
	}

	for _, topup := range topupOrderHistoryResult {
		EMoneyTopUpOrderHistory = append(EMoneyTopUpOrderHistory, user.UserEMoneyTransactionHistory{
			Transaction_amount: topup.Emoney_amont,
			Transaction_type:   "Top Up",
			Transaction_date:   topup.Emoney_updated_at,
		})
	}

	EMoneyTransactionHistory := append(EMoneyOrderHistory, EMoneyOrderSellerHistory...)
	EMoneyTransactionHistory = append(EMoneyTransactionHistory, EMoneyTopUpOrderHistory...)

	layout := "2006-01-02 15:04:05"
	sort.Slice(EMoneyTransactionHistory, func(i, j int) bool {
		date1, _ := time.Parse(layout, EMoneyTransactionHistory[i].Transaction_date)
		date2, _ := time.Parse(layout, EMoneyTransactionHistory[j].Transaction_date)
		return date1.Before(date2)
	})

	for i := range EMoneyTransactionHistory {
		log.Println(EMoneyTransactionHistory[i])
	}

	return EMoneyTransactionHistory
}

func (Usecase *UserUsecase) CheckUserStatus(ctx context.Context, uuid string, costumeid int) user.CheckUserStatusResponse {
	log.Printf("User with uuid: %s enter User Usecase: CheckUserStatus", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = Usecase.CostumeRepository.CheckOwnership(ctx, tx, uuid, costumeid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	statusResult, err := Usecase.UserRepository.CheckUserStatus(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	statusResult.Status = "true"

	return statusResult
}

func (Usecase *UserUsecase) GetSellerAddressDetailByCostumeId(ctx context.Context, userUUID string, costumeID int) user.SellerAddressResponse {
	log.Printf("User with uuid: %s enter User Usecase: GetSellerAddressDetailByCostumeId", userUUID)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, userUUID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	sellerResult, err := Usecase.CostumeRepository.GetSellerIdFindByCostumeID(ctx, tx, userUUID, costumeID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	sellerAddressResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, sellerResult)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	sellerAddressResponse := user.SellerAddressResponse{
		Seller_name:                 sellerAddressResult.Name,
		Seller_origin_province_name: sellerAddressResult.Origin_province_name,
		Seller_origin_province_id:   sellerAddressResult.Origin_city_id,
		Seller_origin_city_name:     sellerAddressResult.Origin_city_name,
		Seller_origin_city_id:       sellerAddressResult.Origin_city_id,
	}

	return sellerAddressResponse
}

func (Usecase *UserUsecase) CheckSellerStatus(ctx context.Context, uuid string) user.CheckUserStatusResponse {
	log.Printf("User with uuid: %s enter User Usecase: CheckSellerStatus", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	statusResult, err := Usecase.UserRepository.CheckUserStatus(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	statusResult.Status = "true"

	return statusResult
}
