package ports

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
)

type Service interface {
	UpdateRanking(ctx context.Context, dto domain.VideoRankingUpdate) error
	GetTopRanked(ctx context.Context, page, limit int) ([]domain.VideoDetail, error)
}
