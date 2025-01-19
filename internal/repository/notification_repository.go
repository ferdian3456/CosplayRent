package repository

import (
	"context"
	"cosplayrent/internal/model/domain"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
)

type NotificationRepository struct {
	Log *zerolog.Logger
}

func NewNotificationRepository(zerolog *zerolog.Logger) *NotificationRepository {
	return &NotificationRepository{
		Log: zerolog,
	}
}

func (repository *NotificationRepository) FindNotificationTemplateById(ctx context.Context, tx *sql.Tx, templateid int) (domain.Notification, error) {
	query := "SELECT template_subject,template_body FROM notifications WHERE id=$1"
	row, err := tx.QueryContext(ctx, query, templateid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	notification := domain.Notification{}

	if row.Next() {
		err = row.Scan(&notification.Template_subject, &notification.Template_body)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}
		return notification, nil
	} else {
		return notification, errors.New("notification template not found")
	}
}
