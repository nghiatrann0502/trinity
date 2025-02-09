package repositories

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"github.com/nghiatrann0502/trinity/proto/gen/proto"
)

type grpcVideoRepository struct {
	log    logger.Logger
	client proto.VideoServiceClient
}

var _ ports.VideoRepository = (*grpcVideoRepository)(nil)

func NewGrpcVideoRepository(log logger.Logger, client proto.VideoServiceClient) ports.VideoRepository {
	return &grpcVideoRepository{
		log:    log.With(map[string]interface{}{"component": "grpc_video_repository"}),
		client: client,
	}
}

func (r *grpcVideoRepository) GetByID(ctx context.Context, id int) (*domain.VideoDetail, error) {
	protoVideo, err := r.client.GetByID(ctx, &proto.VideoRequest{Id: int64(id)})
	if err != nil {
		return nil, err
	}

	if protoVideo == nil || protoVideo.Video == nil {
		return nil, nil
	}

	return &domain.VideoDetail{
		ID:          int(protoVideo.Video.Id),
		Description: protoVideo.Video.Title,
		Title:       protoVideo.Video.Description,
		Thumbnail:   protoVideo.Video.Thumbnail,
		Duration:    int(protoVideo.Video.Duration),
		URL:         protoVideo.Video.Url,
	}, nil
}

func (r *grpcVideoRepository) GetByIDs(ctx context.Context, ids []int) ([]domain.VideoDetail, error) {
	protoIds := make([]int64, len(ids))
	for i, id := range ids {
		protoIds[i] = int64(id)
	}
	protoVideos, err := r.client.GetByIDs(ctx, &proto.GetByIDsRequest{Ids: protoIds})
	if err != nil {
		return nil, err
	}

	videos := make([]domain.VideoDetail, len(protoVideos.Videos))
	for i, video := range protoVideos.Videos {
		videos[i] = domain.VideoDetail{
			ID:          int(video.Id),
			Title:       video.Title,
			Description: video.Description,
			Thumbnail:   video.Thumbnail,
			Duration:    int(video.Duration),
			URL:         video.Url,
		}
	}

	return videos, nil
}
