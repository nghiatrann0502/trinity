package ports

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/video/core/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (*domain.Video, error)
	GetByIDs(ctx context.Context, id []int) ([]domain.Video, error)
}
