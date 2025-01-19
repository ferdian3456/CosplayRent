package usecase

import (
	"bytes"
	"context"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
	"html/template"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587

type NotificationUsecase struct {
	NotificationRepository *repository.NotificationRepository
	DB                     *sql.DB
	Validate               *validator.Validate
	Log                    *zerolog.Logger
	Config                 *koanf.Koanf
}

func NewNotificationUsecase(notificationRepository *repository.NotificationRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *NotificationUsecase {
	return &NotificationUsecase{
		NotificationRepository: notificationRepository,
		DB:                     DB,
		Validate:               validate,
		Log:                    zerolog,
		Config:                 koanf,
	}
}

func (usecase *NotificationUsecase) SendRegisterNotification(ctx context.Context, tx *sql.Tx, username string, useremail string, code string) {
	notification, err := usecase.NotificationRepository.FindNotificationTemplateById(ctx, tx, 1)
	if err != nil {
		respErr := errors.New("failed to find notification template")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	template, err := template.New("emailtemplate").Parse(notification.Template_body)
	if err != nil {
		respErr := errors.New("failed to parse html template")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	data := domain.EmailNotification{
		Username: username,
		Code:     code,
	}

	var tmpl bytes.Buffer
	err = template.Execute(&tmpl, data)
	if err != nil {
		respErr := errors.New("failed to execute html template")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	CONFIG_SENDER_NAME := usecase.Config.String("CONFIG_SENDER_NAME")
	CONFIG_AUTH_EMAIL := usecase.Config.String("CONFIG_AUTH_EMAIL")
	CONFIG_AUTH_PASSWORD := usecase.Config.String("CONFIG_AUTH_PASSWORD")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", useremail)
	mailer.SetHeader("Subject", notification.Template_subject)
	mailer.SetBody("text/html", tmpl.String())

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		respErr := errors.New("failed to send register notification")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}
}
