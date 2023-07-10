package namesdb

import (
	"context"
	"errors"

	"github.com/mwasilew2/go-service-template/internal/domain/models"
)

type NamesDB struct {
	database map[int64]string
}

func NewNamesDB() *NamesDB {
	return &NamesDB{
		database: map[int64]string{
			1: "John",
			2: "Jane",
		},
	}
}

var ErrNameNotFound = errors.New("name not found")

func (n NamesDB) GetName(ctx context.Context, id int64) (*models.Name, error) {
	name, ok := n.database[id]
	if !ok {
		return nil, ErrNameNotFound
	}
	return &models.Name{
		Id:    id,
		Value: name,
	}, nil
}
