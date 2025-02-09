package repositories

import (
	"context"
	"strconv"

	"github.com/nghiatrann0502/trinity/internal/ranking/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/common"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"github.com/nghiatrann0502/trinity/pkg/redisc"
	"github.com/redis/go-redis/v9"
)

type redisCacheRepository struct {
	log    logger.Logger
	client redisc.RedisEngine
}

// GetTopRankedVideos implements ports.CacheRepository.
func (r *redisCacheRepository) GetTopRankedVideos(ctx context.Context, page int, limit int) ([]int, error) {
	videoIDs, err := r.client.GetRedis().ZRevRange(ctx, "ranking", int64((page-1)*limit), (int64(page)*int64(limit))-1).Result()
	if err != nil {
		return nil, common.NewInternalError(err)
	}

	results := make([]int, len(videoIDs))
	for i, id := range videoIDs {
		videoID, err := strconv.Atoi(id)
		if err != nil {
			r.log.Error("failed to convert video id to int", err, nil)
			continue
		}

		results[i] = videoID
	}

	return results, nil
}

var _ ports.CacheRepository = (*redisCacheRepository)(nil)

// NewRedisCacheRepository creates a new instance of redisCacheRepository.
func NewRedisCacheRepository(log logger.Logger, client redisc.RedisEngine) ports.CacheRepository {
	return &redisCacheRepository{
		log:    log.With(map[string]interface{}{"component": "redis_cache_repository"}),
		client: client,
	}
}

// CheckExistKey implements ports.CacheRepository.
func (r *redisCacheRepository) CheckExistKey(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.GetRedis().Exists(ctx, "ranking").Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// CreateNewRankingKey implements ports.CacheRepository.
func (r *redisCacheRepository) CreateNewRankingKey(ctx context.Context, videoID int, value float64) error {
	_, err := r.client.GetRedis().ZAdd(ctx, "ranking", redis.Z{Score: value, Member: videoID}).Result()
	if err != nil {
		return err
	}

	return nil
}

// IncreaseVideoScore implements ports.CacheRepository.
func (r *redisCacheRepository) IncreaseVideoScore(ctx context.Context, videoID int, value float64) error {
	_, err := r.client.GetRedis().ZIncrBy(ctx, "ranking", value, strconv.Itoa(videoID)).Result()
	if err != nil {
		return err
	}
	return nil
}
