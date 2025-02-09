package ports

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
)

type VideoRepository interface {
	GetByID(ctx context.Context, id int) (*domain.VideoDetail, error)
	GetByIDs(ctx context.Context, ids []int) ([]domain.VideoDetail, error)
}
