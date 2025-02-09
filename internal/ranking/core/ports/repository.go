package ports

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
)

type Repository interface {
	GetVideoRankingByVideoId(ctx context.Context, videoId int) (*domain.VideoRanking, error)
	CreateVideoInteraction(ctx context.Context, dto domain.VideoRankingUpdate) error
}
