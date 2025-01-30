package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/repository"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/knadh/koanf/v2"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	googleuuid "github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepository      *repository.UserRepository
	CostumeRepository   *repository.CostumeRepository
	NotificationUsecase *NotificationUsecase
	DB                  *sql.DB
	Validate            *validator.Validate
	Log                 *zerolog.Logger
	Config              *koanf.Koanf
}

func NewUserUsecase(userRepository *repository.UserRepository, costumeRepository *repository.CostumeRepository, notificationUsecase *NotificationUsecase, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *UserUsecase {
	return &UserUsecase{
		UserRepository:      userRepository,
		CostumeRepository:   costumeRepository,
		NotificationUsecase: notificationUsecase,
		DB:                  DB,
		Validate:            validate,
		Log:                 zerolog,
		Config:              koanf,
	}
}

func (usecase *UserUsecase) Create(ctx context.Context, request user.UserCreateRequest) (string, error) {
	err := usecase.Validate.Struct(request)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return "", respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		respErr := errors.New("error generating password hash")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	uuid := googleuuid.New()

	user := domain.User{
		Id:         uuid.String(),
		Name:       request.Name,
		Email:      request.Email,
		Password:   string(hashedPassword),
		Created_at: &now,
	}

	err = usecase.UserRepository.CheckCredentialUnique(ctx, tx, user)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return "", err
	}

	usecase.UserRepository.Create(ctx, tx, user)

	secretKey := usecase.Config.String("SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.Id,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		respErr := errors.New("failed to sign a token")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 5)
	for i := range code {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[randomIndex.Int64()]
	}

	expiredAt := now.Add(5 * time.Minute)

	UserVerification := domain.UserVerification{
		User_id:           user.Id,
		Verification_code: string(code),
		Created_at:        &now,
		Updated_at:        &now,
		Expired_at:        &expiredAt,
	}

	usecase.UserRepository.CreateUserVerification(ctx, tx, UserVerification)

	usecase.NotificationUsecase.SendRegisterNotification(ctx, tx, user.Name, user.Email, string(code))

	return tokenString, nil
}

func (usecase *UserUsecase) Login(ctx context.Context, request user.UserLoginRequest) (string, error) {
	err := usecase.Validate.Struct(request)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return "", respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	userRequest := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user := domain.User{}
	user, err = usecase.UserRepository.Login(ctx, tx, userRequest.Email)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		respErr := errors.New("wrong password")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return "", respErr
	}

	secretKey := usecase.Config.String("SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.Id,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		respErr := errors.New("failed to sign a token")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	return tokenString, nil
}

func (usecase *UserUsecase) FindByUUID(ctx context.Context, uuid string) (user.UserResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := user.UserResponse{}
	user, err = usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	if user.Profile_picture != nil {
		value := imageEnv + *user.Profile_picture
		user.Profile_picture = &value
	}

	return user, nil
}

func (usecase *UserUsecase) CheckUserExistance(ctx context.Context, uuid string) error {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	err = usecase.UserRepository.CheckUserExistance(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	return nil
}

func (usecase *UserUsecase) CheckUserExistanceForNonActivated(ctx context.Context, uuid string) error {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	err = usecase.UserRepository.CheckUserExistanceForNonActivated(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	return nil
}

func (usecase *UserUsecase) VerifyCode(ctx context.Context, request user.UserVerificationCode, uuid string) error {
	err := usecase.Validate.Struct(request)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user, err := usecase.UserRepository.VerifyCode(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	if request.Code == user.Verification_code {
		now := time.Now()
		fmt.Println("now:", &now)
		fmt.Println("verification code expired_at:", user.Expired_at)
		if now.After(*user.Expired_at) {
			respErr := errors.New("verification code expired")
			usecase.Log.Warn().Err(respErr).Msg(respErr.Error())
			return respErr
		}
		usecase.UserRepository.ChangeVerificationStatus(ctx, tx, uuid)
		return nil
	} else {
		respErr := errors.New("invalid verification code")
		usecase.Log.Warn().Err(respErr).Msg(respErr.Error())
		return respErr
	}
}

func (usecase *UserUsecase) FindAll(ctx context.Context, uuid string) ([]user.UserResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := []user.UserResponse{}

	user, err = usecase.UserRepository.FindAll(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range user {
		if user[i].Profile_picture != nil {
			value := imageEnv + *user[i].Profile_picture
			user[i].Profile_picture = &value
		}
	}

	return user, nil
}

func (usecase *UserUsecase) Update(ctx context.Context, userRequest user.UserPatchRequest, uuid string) error {
	err := usecase.Validate.Struct(userRequest)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()

	if userRequest.Profile_picture != nil {
		user := domain.User{
			Id:                   uuid,
			Name:                 *userRequest.Name,
			Email:                *userRequest.Email,
			Address:              *userRequest.Address,
			Profile_picture:      userRequest.Profile_picture,
			Origin_province_name: *userRequest.Origin_province_name,
			Origin_province_id:   *userRequest.Origin_province_id,
			Origin_city_name:     *userRequest.Origin_city_name,
			Origin_city_id:       *userRequest.Origin_city_id,
			Updated_at:           &now,
		}

		usecase.UserRepository.Update(ctx, tx, user)

		return nil
	} else {
		user := domain.User{
			Id:                   uuid,
			Name:                 *userRequest.Name,
			Email:                *userRequest.Email,
			Address:              *userRequest.Address,
			Origin_province_name: *userRequest.Origin_province_name,
			Origin_province_id:   *userRequest.Origin_province_id,
			Origin_city_name:     *userRequest.Origin_city_name,
			Origin_city_id:       *userRequest.Origin_city_id,
			Updated_at:           &now,
		}

		usecase.UserRepository.Update(ctx, tx, user)

		return nil
	}
}

func (usecase *UserUsecase) Delete(ctx context.Context, uuid string) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	userProfile, err := usecase.UserRepository.FindProfileById(ctx, tx, uuid)
	if err != nil {
		respErr := errors.New("profile picture not found")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
	}

	if userProfile != nil {
		finalProfilePicturePath := ".." + *userProfile
		err = os.Remove(finalProfilePicturePath)
		if err != nil {
			respErr := errors.New("failed to remove profile picture")
			usecase.Log.Warn().Err(respErr).Msg(err.Error())
		}
	}

	usecase.UserRepository.Delete(ctx, tx, uuid)
}

func (usecase *UserUsecase) AddOrUpdateIdentityCard(ctx context.Context, uuid string, userRequest user.IdentityCardRequest) error {
	err := usecase.Validate.Struct(userRequest)
	if err != nil {
		respErr := errors.New("identity card cannot be empty")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Id:                    uuid,
		Identity_card_picture: *userRequest.IdentityCard_picture,
	}

	usecase.UserRepository.AddOrUpdateIdentityCard(ctx, tx, user)

	return nil
}

func (usecase *UserUsecase) GetIdentityCard(ctx context.Context, uuid string) (user.IdentityCardResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := user.IdentityCardResponse{}
	identityCardResult, err := usecase.UserRepository.GetIdentityCard(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	user.IdentityCard_picture = imageEnv + identityCardResult
	return user, nil
}

func (usecase *UserUsecase) GetEMoneyAmount(ctx context.Context, uuid string) user.UserEmoneyResponse {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := usecase.UserRepository.GetEMoneyAmount(ctx, tx, uuid)

	return user
}

func (usecase *UserUsecase) GetEMoneyTransactionHistory(ctx context.Context, uuid string) ([]user.UserEMoneyTransactionHistory, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user, err := usecase.UserRepository.FindAllMoneyChanges(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	return user, nil
}

func (usecase *UserUsecase) CheckUserStatus(ctx context.Context, uuid string, costumeid int) (user.CheckUserStatusResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := user.CheckUserStatusResponse{}

	err = usecase.CostumeRepository.CheckOwnership(ctx, tx, uuid, costumeid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	user, err = usecase.UserRepository.CheckUserStatus(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	user.Status = "true"

	return user, nil
}

func (usecase *UserUsecase) CheckSellerStatus(ctx context.Context, uuid string) (user.CheckUserStatusResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := user.CheckUserStatusResponse{}

	user, err = usecase.UserRepository.CheckUserStatus(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	user.Status = "true"

	return user, nil
}

func (usecase *UserUsecase) FindSellerAddressDetailByCostumeId(ctx context.Context, userUUID string, costumeID int) (user.SellerAddressResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	sellerResult, err := usecase.CostumeRepository.FindSellerIdFindByCostumeID(ctx, tx, costumeID)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user.SellerAddressResponse{}, err
	}

	sellerAddressResult, err := usecase.UserRepository.FindAddressByUserId(ctx, tx, sellerResult)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user.SellerAddressResponse{}, err
	}

	sellerAddressResponse := user.SellerAddressResponse{
		Seller_name:                 sellerAddressResult.Name,
		Seller_origin_province_name: &sellerAddressResult.Origin_province_name,
		Seller_origin_province_id:   &sellerAddressResult.Origin_city_id,
		Seller_origin_city_name:     &sellerAddressResult.Origin_city_name,
		Seller_origin_city_id:       &sellerAddressResult.Origin_city_id,
	}

	return sellerAddressResponse, nil
}

//func (usecase *UserUsecase) CheckSellerStatus(ctx context.Context, uuid string) user.CheckUserStatusResponse {
//	log.Printf("User with uuid: %s enter User Usecase: CheckSellerStatus", uuid)
//
//	tx, err := usecaseDB.Begin()
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	defer helper.CommitOrRollback(tx)
//
//	_, err = usecaseUserRepository.FindByUUID(ctx, tx, uuid)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	statusResult, err := usecaseUserRepository.CheckUserStatus(ctx, tx, uuid)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	statusResult.Status = "true"
//
//	return statusResult
//}
