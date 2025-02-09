package services

import (
	"context"
	"time"

	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/common"
	"github.com/nghiatrann0502/trinity/pkg/logger"
)

type service struct {
	repository      ports.Repository
	cache           ports.CacheRepository
	videoRepository ports.VideoRepository
	log             logger.Logger
}

var _ ports.Service = (*service)(nil)

func NewRankingService(log logger.Logger, repository ports.Repository, cache ports.CacheRepository, videoRPC ports.VideoRepository) ports.Service {
	return &service{
		log:             log.With(map[string]interface{}{"component": "ranking_service"}),
		repository:      repository,
		cache:           cache,
		videoRepository: videoRPC,
	}
}

// Fake calculate ranking
// TODO: Implement real ranking calculation
func (s *service) calculateRanking(action string, value float64) float64 {
	switch action {
	case "view":
		return float64(domain.ViewWeight) * value
	case "like":
		return float64(domain.LikeWeight) * value
	case "comment":
		return float64(domain.CommentWeight) * value
	case "share":
		return float64(domain.ShareWeight) * value
	case "watch_time":
		return domain.WatchWeight * value
	default:
		return 0
	}
}

func (s *service) UpdateRanking(ctx context.Context, dto domain.VideoRankingUpdate) error {
	// Validate video ranking update
	if ok := domain.ValidateAction(dto.Action); !ok {
		return common.NewValidationError("Invalid action")
	}

	video, err := s.videoRepository.GetByID(ctx, dto.VideoID)
	if err != nil {
		return common.NewInternalError(err)
	}

	if video == nil {
		return common.NewNotFoundError("Video not found")
	}

	// Create video interaction
	// Retry with backoff if error, you can use retry package from github.com/nghiatrann0502/trinity/pkg/common
	// Use goroutine to run this function in background to avoid blocking, bottleneck when create video interaction and fast response to client
	go common.RetryWithBackoff(5, time.Second, func() error {
		return s.repository.CreateVideoInteraction(ctx, dto)
	})

	// Get video ranking by video id
	// videoRanking, err := s.repository.GetVideoRankingByVideoId(ctx, dto.VideoID)
	// if err != nil {
	// 	return nil
	// }

	// Calculate increment score
	increment := s.calculateRanking(dto.Action, float64(dto.Value))

	exist, err := s.cache.CheckExistKey(ctx, "ranking")
	if err != nil {
		return common.NewInternalError(err)
	}

	// Check if ranking key exist
	if !exist {
		// Create new ranking
		// cClient.ZAdd(ctx, "ranking", redis.Z{Score: increment, Member: dto.VideoID}).Result()
		if err := s.cache.CreateNewRankingKey(ctx, dto.VideoID, increment); err != nil {
			return common.NewInternalError(err)
		}
	} else {
		// _, err := cClient.ZIncrBy(ctx, "ranking", increment, strconv.Itoa(dto.VideoID)).Result()
		// if err != nil {
		// 	return common.NewInternalError(err)
		// }

		// Increase video score
		if err := s.cache.IncreaseVideoScore(ctx, dto.VideoID, increment); err != nil {
			return common.NewInternalError(err)
		}
	}

	// NOTE: Why we don't update database here?
	// Because we want to update database in worker to avoid bottleneck when update ranking
	// We can use message queue to send message to worker to update database, cronjob or other way to update database
	// Best practice is using message queue (kafka for high performance) to send message to worker to update database (I think so!!!)

	return nil
}

func (s *service) GetTopRanked(ctx context.Context, page, limit int) ([]domain.VideoDetail, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	videoIDs, err := s.cache.GetTopRankedVideos(ctx, page, limit)
	if err != nil {
		return nil, common.NewInternalError(err)
	}

	videoDetails, err := s.videoRepository.GetByIDs(ctx, videoIDs)
	if err != nil {
		return nil, common.NewInternalError(err)
	}

	// NOTE: Reduce time complexity from O(n^2) to O(n)
	// But space complexity increase to O(n)
	mDetails := make(map[int]domain.VideoDetail)
	for _, v := range videoDetails {
		mDetails[v.ID] = v
	}

	// TODO: Get video detail by video ids using grpc client or gpc api
	var videos []domain.VideoDetail

	for _, id := range videoIDs {
		if v, ok := mDetails[id]; ok {
			videos = append(videos, v)
		}
	}

	return videos, nil
}
