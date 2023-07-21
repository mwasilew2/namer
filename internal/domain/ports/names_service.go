package ports

import (
	"context"

	"github.com/mwasilew2/go-service-template/internal/domain/models"
)

type NamesService interface {
	GetName(ctx context.Context, year int64, id int64) (*models.Name, error)
	GetPage(ctx context.Context, year int64, page int64, limit int64) ([]*models.Name, error)
	GetYearsAvailable(ctx context.Context) (map[int64]struct{}, error)
	GetNoOfEntries(ctx context.Context, year int64) (int64, error)
}
