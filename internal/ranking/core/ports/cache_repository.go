package ports

import "context"

type CacheRepository interface {
	CheckExistKey(ctx context.Context, key string) (bool, error)
	CreateNewRankingKey(ctx context.Context, videoID int, value float64) error
	IncreaseVideoScore(ctx context.Context, videoID int, value float64) error
	GetTopRankedVideos(ctx context.Context, page, limit int) ([]int, error)
}
