package ports

import (
	"context"

	"github.com/mwasilew2/go-service-template/internal/domain/models"
)

type NamesService interface {
	GetName(ctx context.Context, id int64) (*models.Name, error)
}
