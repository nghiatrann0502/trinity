package ports

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/video/core/domain"
)

type Service interface {
	GetVideoByID(ctx context.Context, id int) (*domain.Video, error)
	GetVideoByIDs(ctx context.Context, id []int) ([]domain.Video, error)
}
