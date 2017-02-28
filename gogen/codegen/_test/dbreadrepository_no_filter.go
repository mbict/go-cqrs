package sql

import (
	"database/sql"
	"github.com/masterminds/squirrel"
	"github.com/satori/go.uuid"
	"testing/base/models"
	"testing/base/repository"
)

type dbProductitemRepository struct {
	db *sql.DB
}

var (
	productitemQuery = squirrel.Select("id", "name", "price").From("productitem")
)

func NewProductitemRepository(db *sql.DB) repository.ProductitemRepository {
	return &dbProductitemRepository{
		db: db,
	}
}

func (r *dbProductitemRepository) FindAll() ([]*models.Item, error) {
	rows, err := productitemQuery.RunWith(r.db).Query()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var result []*models.Item
	for rows.Next() {
		item := models.Item{}
		err = rows.Scan(&item.Id, &item.Name, &item.Price)
		if err != nil {
			return nil, err
		}
		result = append(result, &item)
	}
	return result, nil
}

func (r *dbProductitemRepository) Find(id uuid.UUID) (*models.Item, error) {
	var item models.Item
	row := productitemQuery.Where("id = ?", id).RunWith(r.db).QueryRow()
	err := row.Scan(&item.Id, &item.Name, &item.Price)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &item, nil
}