package repositories

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/video/core/domain"
	"github.com/nghiatrann0502/trinity/internal/video/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/database"
	"github.com/nghiatrann0502/trinity/pkg/logger"
)

type mysqlRepo struct {
	log logger.Logger
	db  database.DBEngine
}

var _ ports.Repository = (*mysqlRepo)(nil)

func NewMysqlRepo(log logger.Logger, db database.DBEngine) ports.Repository {
	return &mysqlRepo{
		log: log.With(map[string]interface{}{"component": "mysql-repo"}),
		db:  db,
	}
}

func (r *mysqlRepo) GetByID(ctx context.Context, id int) (*domain.Video, error) {
	r.log.Info("GetByID", nil)
	return nil, nil
}

func (r *mysqlRepo) GetByIDs(ctx context.Context, ids []int) ([]domain.Video, error) {
	r.log.Info("GetByIDs", nil)
	return nil, nil
}
