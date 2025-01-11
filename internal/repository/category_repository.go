package repository

import (
	"context"
	"cosplayrent/internal/model/web/category"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
)

type CategoryRepository struct {
	Log *zerolog.Logger
}

func NewCategoryRepository(zerolog *zerolog.Logger) *CategoryRepository {
	return &CategoryRepository{
		Log: zerolog,
	}
}

func (repository *CategoryRepository) FindAllCategory(ctx context.Context, tx *sql.Tx) ([]category.CategoryResponse, error) {
	query := "SELECT id,name from categories"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	hasData := false

	defer rows.Close()

	categories := []category.CategoryResponse{}
	for rows.Next() {
		category := category.CategoryResponse{}
		err = rows.Scan(&category.Id, &category.Name)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}

		categories = append(categories, category)
		hasData = true
	}

	if hasData == false {
		return categories, errors.New("category is empty")
	}

	return categories, nil
}

func (repository *CategoryRepository) FindCategoryNameById(ctx context.Context, tx *sql.Tx, id int) (string, error) {
	query := "SELECT name from categories where id=$1"
	row, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer row.Close()

	var categoryName string

	if row.Next() {
		err := row.Scan(&categoryName)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}

		return categoryName, nil
	} else {
		return categoryName, errors.New("category not found")
	}
}
