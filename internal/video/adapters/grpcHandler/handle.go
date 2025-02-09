package grpchandler

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/video/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"github.com/nghiatrann0502/trinity/proto/gen/proto"
)

type gRPCHandler struct {
	proto.UnimplementedVideoServiceServer

	log     logger.Logger
	service ports.Service
}

var _ proto.VideoServiceServer = (*gRPCHandler)(nil)

func NewGRPCHandler(log logger.Logger, service ports.Service) proto.VideoServiceServer {
	return &gRPCHandler{
		log:     log.With(map[string]interface{}{"component": "grpc-handler"}),
		service: service,
	}
}

func (h *gRPCHandler) GetByID(ctx context.Context, req *proto.VideoRequest) (*proto.VideoResponse, error) {
	video, err := h.service.GetVideoByID(ctx, int(req.Id))
	if err != nil {
		h.log.Error("failed to get video by id", err, nil)
		return nil, err
	}

	if video == nil {
		return nil, nil
	}

	return &proto.VideoResponse{
		Video: &proto.VideoDetail{
			Id:          int64(video.ID),
			Title:       video.Title,
			Description: video.Description,
			Thumbnail:   video.Thumbnail,
			Url:         video.URL,
			Duration:    int64(video.Duration),
		},
	}, nil
}

func (h *gRPCHandler) GetByIDs(ctx context.Context, req *proto.GetByIDsRequest) (*proto.VideoList, error) {
	ids := make([]int, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = int(id)
	}

	videos, err := h.service.GetVideoByIDs(ctx, ids)
	if err != nil {
		h.log.Error("failed to get videos by ids", err, nil)
		return nil, err
	}

	protoVideos := make([]*proto.VideoDetail, len(videos))
	for i, video := range videos {
		protoVideos[i] = &proto.VideoDetail{
			Id:          int64(video.ID),
			Title:       video.Title,
			Description: video.Description,
			Thumbnail:   video.Thumbnail,
			Url:         video.URL,
			Duration:    int64(video.Duration),
		}
	}

	return &proto.VideoList{
		Videos: protoVideos,
	}, nil
}
