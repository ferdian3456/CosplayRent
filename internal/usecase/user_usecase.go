package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"github.com/knadh/koanf/v2"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	googleuuid "github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepository *repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
	Log            *zerolog.Logger
	Config         *koanf.Koanf
}

func NewUserUsecase(userRepository *repository.UserRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *UserUsecase {
	return &UserUsecase{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
		Log:            zerolog,
		Config:         koanf,
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
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		respErr := errors.New("error generating password hash")
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
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

	secretKey := usecase.Config.String("application.secret_key")
	secretKeyByte := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.Id,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		respErr := errors.New("failed to sign a token")
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
	}

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
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
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

	secretKey := usecase.Config.String("application.secret_key")
	secretKeyByte := []byte(secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.Id,
		"expired": time.Date(2030, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		respErr := errors.New("failed to sign a token")
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
	}

	return tokenString, nil
}

func (usecase *UserUsecase) FindByUUID(ctx context.Context, uuid string) (user.UserResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := user.UserResponse{}
	user, err = usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	imageEnv := usecase.Config.String("application.image_env")

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
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer helper.CommitOrRollback(tx)

	err = usecase.UserRepository.CheckUserExistance(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	return nil
}

func (usecase *UserUsecase) FindAll(ctx context.Context, uuid string) ([]user.UserResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer helper.CommitOrRollback(tx)

	user := []user.UserResponse{}

	user, err = usecase.UserRepository.FindAll(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return user, err
	}

	imageEnv := usecase.Config.String("application.image_env")

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
		usecase.Log.Panic().Err(respErr).Msg(err.Error())
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()

	user := domain.User{
		Id:                   uuid,
		Name:                 *userRequest.Name,
		Email:                *userRequest.Email,
		Address:              *userRequest.Address,
		Profile_picture:      *userRequest.Profile_picture,
		Origin_province_name: *userRequest.Origin_province_name,
		Origin_province_id:   *userRequest.Origin_province_id,
		Origin_city_name:     *userRequest.Origin_city_name,
		Origin_city_id:       *userRequest.Origin_city_id,
		Updated_at:           &now,
	}

	usecase.UserRepository.Update(ctx, tx, user)

	return nil
}

//
//func (usecase *UserUsecase) Delete(ctx context.Context, uuid string) {
//	log.Printf("User with uuid: %s enter User Usecase: Delete", uuid)
//
//	tx, err := usecaseDB.Begin()
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	defer helper.CommitOrRollback(tx)
//
//	userResult, err := usecaseUserRepository.FindByUUID(ctx, tx, uuid)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	usecaseUserRepository.Delete(ctx, tx, uuid)
//
//	finalProfilePicturePath := ".." + *userResult.Profile_picture
//
//	err = os.Remove(finalProfilePicturePath)
//	helper.PanicIfError(err)
//}
//
//func (usecase *UserUsecase) VerifyAndRetrieve(ctx context.Context, tokenString string) (user.UserResponse, error) {
//	secretKey := os.Getenv("SECRET_KEY")
//	secretKeyByte := []byte(secretKey)
//
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//		return secretKeyByte, nil
//	})
//
//	if err != nil || !token.Valid {
//		return user.UserResponse{}, errors.New("token is not valid")
//	}
//
//	var id string
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		if val, exists := claims["id"]; exists {
//			if strVal, ok := val.(string); ok {
//				id = strVal
//			} else {
//				return user.UserResponse{}, fmt.Errorf("id claim is not a string")
//			}
//		} else {
//			return user.UserResponse{}, fmt.Errorf("id claim does not exist")
//		}
//	}
//
//	tx, err := usecaseDB.Begin()
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	defer helper.CommitOrRollback(tx)
//	userDomain, err := usecaseUserRepository.FindByUUID(ctx, tx, id)
//	helper.PanicIfError(err)
//
//	return userDomain, nil
//}
//
//func (usecase *UserUsecase) AddIdentityCard(ctx context.Context, uuid string, IdentityCardImage string) {
//	log.Printf("User with uuid: %s enter User Usecase: AddIdentityCard", uuid)
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
//	usecaseUserRepository.AddOrUpdateIdentityCard(ctx, tx, uuid, IdentityCardImage)
//}
//
//func (usecase *UserUsecase) GetIdentityCard(ctx context.Context, uuid string) (identityCardImage string) {
//	log.Printf("User with uuid: %s enter User Usecase: GetIdentityCard", uuid)
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
//	identityCardResult, err := usecaseUserRepository.GetIdentityCard(ctx, tx, uuid)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	err = godotenv.Load("../.env")
//	helper.PanicIfError(err)
//
//	imageEnv := os.Getenv("IMAGE_ENV")
//
//	identityCardResult = imageEnv + identityCardResult
//
//	return identityCardResult
//}
//
//func (usecase *UserUsecase) UpdateIdentityCard(ctx context.Context, uuid string, IdentityCardImage string) {
//	log.Printf("User with uuid: %s enter User Usecase: UpdateIdentityCard", uuid)
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
//	usecaseUserRepository.AddOrUpdateIdentityCard(ctx, tx, uuid, IdentityCardImage)
//}
//
//func (usecase *UserUsecase) GetEMoneyAmount(ctx context.Context, uuid string) (userEmoneyResult user.UserEmoneyResponse) {
//	log.Printf("User with uuid: %s enter User Usecase: GetEMoneyAmount", uuid)
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
//	userEmoneyResult, err = usecaseUserRepository.GetEMoneyAmount(ctx, tx, uuid)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	return userEmoneyResult
//}
//
//func (usecase *UserUsecase) GetEMoneyTransactionHistory(ctx context.Context, uuid string) []user.UserEMoneyTransactionHistory {
//	log.Printf("User with uuid: %s enter User Usecase: GetEMoneyTransactionHistory", uuid)
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
//	orderHistoryResult, err1 := usecaseOrderRepository.FindOrderHistoryByUserId(ctx, tx, uuid)
//
//	orderHistorySellerResult, err := usecaseOrderRepository.FindOrderHistoryBySellerId(ctx, tx, uuid)
//
//	topupOrderHistoryResult, err2 := usecaseTopUpOrderRepository.FindTopUpOrderHistoryByUserId(ctx, tx, uuid)
//
//	err3 := errors.New("order and topup order is not found")
//
//	if err != nil && err1 != nil && err2 != nil {
//		panic(exception.NewNotFoundError(err3.Error()))
//	}
//
//	EMoneyOrderHistory := []user.UserEMoneyTransactionHistory{}
//	EMoneyOrderSellerHistory := []user.UserEMoneyTransactionHistory{}
//	EMoneyTopUpOrderHistory := []user.UserEMoneyTransactionHistory{}
//
//	for _, order := range orderHistoryResult {
//		EMoneyOrderHistory = append(EMoneyOrderHistory, user.UserEMoneyTransactionHistory{
//			Transaction_amount: order.Emoney_amont,
//			Transaction_type:   "Order (Buyer)",
//			Transaction_date:   order.Emoney_updated_at,
//		})
//	}
//
//	for _, order := range orderHistorySellerResult {
//		EMoneyOrderSellerHistory = append(EMoneyOrderSellerHistory, user.UserEMoneyTransactionHistory{
//			Transaction_amount: order.Emoney_amont,
//			Transaction_type:   "Order (Seller)",
//			Transaction_date:   order.Emoney_updated_at,
//		})
//	}
//
//	for _, topup := range topupOrderHistoryResult {
//		EMoneyTopUpOrderHistory = append(EMoneyTopUpOrderHistory, user.UserEMoneyTransactionHistory{
//			Transaction_amount: topup.Emoney_amont,
//			Transaction_type:   "Top Up",
//			Transaction_date:   topup.Emoney_updated_at,
//		})
//	}
//
//	EMoneyTransactionHistory := append(EMoneyOrderHistory, EMoneyOrderSellerHistory...)
//	EMoneyTransactionHistory = append(EMoneyTransactionHistory, EMoneyTopUpOrderHistory...)
//
//	layout := "2006-01-02 15:04:05"
//	sort.Slice(EMoneyTransactionHistory, func(i, j int) bool {
//		date1, _ := time.Parse(layout, EMoneyTransactionHistory[i].Transaction_date)
//		date2, _ := time.Parse(layout, EMoneyTransactionHistory[j].Transaction_date)
//		return date1.Before(date2)
//	})
//
//	for i := range EMoneyTransactionHistory {
//		log.Println(EMoneyTransactionHistory[i])
//	}
//
//	return EMoneyTransactionHistory
//}
//
//func (usecase *UserUsecase) CheckUserStatus(ctx context.Context, uuid string, costumeid int) user.CheckUserStatusResponse {
//	log.Printf("User with uuid: %s enter User Usecase: CheckUserStatus", uuid)
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
//	err = usecaseCostumeRepository.CheckOwnership(ctx, tx, uuid, costumeid)
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
//
//func (usecase *UserUsecase) GetSellerAddressDetailByCostumeId(ctx context.Context, userUUID string, costumeID int) user.SellerAddressResponse {
//	log.Printf("User with uuid: %s enter User Usecase: GetSellerAddressDetailByCostumeId", userUUID)
//
//	tx, err := usecaseDB.Begin()
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	defer helper.CommitOrRollback(tx)
//
//	_, err = usecaseUserRepository.FindByUUID(ctx, tx, userUUID)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	sellerResult, err := usecaseCostumeRepository.GetSellerIdFindByCostumeID(ctx, tx, userUUID, costumeID)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	sellerAddressResult, err := usecaseUserRepository.FindByUUID(ctx, tx, sellerResult)
//	if err != nil {
//		panic(exception.NewNotFoundError(err.Error()))
//	}
//
//	sellerAddressResponse := user.SellerAddressResponse{
//		Seller_name:                 sellerAddressResult.Name,
//		Seller_origin_province_name: sellerAddressResult.Origin_province_name,
//		Seller_origin_province_id:   sellerAddressResult.Origin_city_id,
//		Seller_origin_city_name:     sellerAddressResult.Origin_city_name,
//		Seller_origin_city_id:       sellerAddressResult.Origin_city_id,
//	}
//
//	return sellerAddressResponse
//}
//
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
