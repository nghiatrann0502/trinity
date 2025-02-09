package services

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/video/core/domain"
	"github.com/nghiatrann0502/trinity/internal/video/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/logger"
)

type service struct {
	log        logger.Logger
	repository ports.Repository
}

var _ ports.Service = (*service)(nil)

func NewService(log logger.Logger, repository ports.Repository) ports.Service {
	return &service{
		log:        log.With(map[string]interface{}{"component": "service"}),
		repository: repository,
	}
}

func (s *service) GetVideoByID(ctx context.Context, id int) (*domain.Video, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *service) GetVideoByIDs(ctx context.Context, ids []int) ([]domain.Video, error) {
	return s.repository.GetByIDs(ctx, ids)
}
