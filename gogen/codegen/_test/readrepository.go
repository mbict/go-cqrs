package repository

import (
	"github.com/satori/go.uuid"
	"testing/base/models"
)


type ItemRepository interface {
	func FindAll() ([]*models.Testmodel1, error)
	func Find(id uuid.UUID) (*models.Testmodel1, error)
}

type ProductFilter struct {
	Id []uuid.UUID
	Name string
	Test []string
}

type ProductRepository interface {
	func FindAll(filter *ProductFilter) ([]*models.Testmodel2, error)
	func Find(id uuid.UUID) (*models.Testmodel2, error)
}
