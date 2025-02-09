package repositories

import (
	"context"

	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/common"
	"github.com/nghiatrann0502/trinity/pkg/database"
	"github.com/nghiatrann0502/trinity/pkg/logger"
)

type mysqlRepository struct {
	log logger.Logger
	db  database.DBEngine
}

var _ ports.Repository = (*mysqlRepository)(nil)

func NewMySQLRepository(log logger.Logger, db database.DBEngine) ports.Repository {
	return &mysqlRepository{
		log: log.With(map[string]interface{}{"component": "mysql_repository"}),
		db:  db,
	}
}

func (mysql *mysqlRepository) GetVideoRankingByVideoId(ctx context.Context, videoId int) (*domain.VideoRanking, error) {
	return nil, nil
}

func (mysql *mysqlRepository) CreateVideoInteraction(ctx context.Context, dto domain.VideoRankingUpdate) error {
	query := `INSERT INTO video_interaction (video_id, action, value) 
  VALUES (?, ?, ?)`

	_, err := mysql.db.GetDB().Exec(query, dto.VideoID, dto.Action, dto.Value)
	if err != nil {
		return common.NewDatabaseError(err)
	}

	return nil
}
