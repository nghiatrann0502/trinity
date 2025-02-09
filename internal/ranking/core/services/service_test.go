package services_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/nghiatrann0502/trinity/internal/ranking/adapters/repositories"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/services"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"go.uber.org/mock/gomock"
)

func TestUpdateRanking(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := repositories.NewMockRepository(mockCtrl)
	mockCache := repositories.NewMockCacheRepository(mockCtrl)
	mockGrpc := repositories.NewMockVideoRepository(mockCtrl)
	mockLogger := logger.NewMockLogger(mockCtrl)

	// Set up logger mock expectation
	mockLogger.EXPECT().With(map[string]interface{}{"component": "ranking_service"}).Return(mockLogger)

	svc := services.NewRankingService(mockLogger, mockRepo, mockCache, mockGrpc)

	testCases := []struct {
		name       string
		dto        domain.VideoRankingUpdate
		setupMocks func()
		wantErr    bool
	}{
		{
			name: "Success - Key exists",
			dto: domain.VideoRankingUpdate{
				VideoID: 1,
				Action:  "view",
				Value:   1,
			},
			setupMocks: func() {
				mockRepo.EXPECT().CreateVideoInteraction(gomock.Any(), domain.VideoRankingUpdate{
					VideoID: 1,
					Action:  "view",
					Value:   1,
				}).Return(nil).AnyTimes()
				mockGrpc.EXPECT().GetByID(gomock.Any(), 1).Return(&domain.VideoDetail{
					ID:          1,
					Title:       "Video title",
					Description: "Video description",
					Thumbnail:   "https://example",
					Duration:    100,
					URL:         "https://example.com",
				}, nil).AnyTimes()
				mockCache.EXPECT().CheckExistKey(gomock.Any(), "ranking").Return(true, nil)
				mockCache.EXPECT().IncreaseVideoScore(gomock.Any(), 1, float64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success - Key does not exist",
			dto: domain.VideoRankingUpdate{
				VideoID: 2,
				Action:  "like",
				Value:   1,
			},
			setupMocks: func() {
				mockRepo.EXPECT().CreateVideoInteraction(gomock.Any(), domain.VideoRankingUpdate{
					VideoID: 2,
					Action:  "like",
					Value:   1,
				}).Return(nil).AnyTimes()
				mockGrpc.EXPECT().GetByID(gomock.Any(), 2).Return(&domain.VideoDetail{
					ID:          2,
					Title:       "Video title",
					Description: "Video description",
					Thumbnail:   "https://example",
					Duration:    100,
					URL:         "https://example.com",
				}, nil).AnyTimes()
				mockCache.EXPECT().CheckExistKey(gomock.Any(), "ranking").Return(false, nil)
				mockCache.EXPECT().CreateNewRankingKey(gomock.Any(), 2, float64(2)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid action",
			dto: domain.VideoRankingUpdate{
				VideoID: 1,
				Action:  "invalid",
				Value:   1,
			},
			setupMocks: func() {},
			wantErr:    true,
		},
		{
			name: "Cache check error",
			dto: domain.VideoRankingUpdate{
				VideoID: 1,
				Action:  "view",
				Value:   1,
			},
			setupMocks: func() {
				mockRepo.EXPECT().CreateVideoInteraction(gomock.Any(), domain.VideoRankingUpdate{
					VideoID: 1,
					Action:  "view",
					Value:   1,
				}).Return(nil).AnyTimes()
				mockCache.EXPECT().CheckExistKey(gomock.Any(), "ranking").Return(false, errors.New("cache error"))
			},
			wantErr: true,
		},
		{
			name: "Cache update error",
			dto: domain.VideoRankingUpdate{
				VideoID: 1,
				Action:  "view",
				Value:   1,
			},
			setupMocks: func() {
				mockRepo.EXPECT().CreateVideoInteraction(gomock.Any(), domain.VideoRankingUpdate{
					VideoID: 1,
					Action:  "view",
					Value:   1,
				}).Return(nil).AnyTimes()
				mockGrpc.EXPECT().GetByID(gomock.Any(), 1).Return(&domain.VideoDetail{
					ID:          1,
					Title:       "Video title",
					Description: "Video description",
					Thumbnail:   "https://example",
					Duration:    100,
					URL:         "https://example.com",
				}, nil).AnyTimes()
				mockCache.EXPECT().CheckExistKey(gomock.Any(), "ranking").Return(true, nil)
				mockCache.EXPECT().IncreaseVideoScore(gomock.Any(), 1, float64(1)).Return(errors.New("update error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := svc.UpdateRanking(context.Background(), tt.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRanking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRanking(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := repositories.NewMockRepository(mockCtrl)
	mockCache := repositories.NewMockCacheRepository(mockCtrl)
	mockLogger := logger.NewMockLogger(mockCtrl)
	mockGrpc := repositories.NewMockVideoRepository(mockCtrl)

	// Set up logger mock expectation
	mockLogger.EXPECT().With(map[string]interface{}{"component": "ranking_service"}).Return(mockLogger)

	svc := services.NewRankingService(mockLogger, mockRepo, mockCache, mockGrpc)

	testCases := []struct {
		name       string
		page       int
		limit      int
		setupMocks func()
		want       []domain.VideoDetail
		wantErr    bool
	}{
		{
			name:  "Success - Empty result",
			page:  1,
			limit: 10,
			setupMocks: func() {
				mockCache.EXPECT().GetTopRankedVideos(gomock.Any(), 1, 10).Return([]int{}, nil)
				mockGrpc.EXPECT().GetByIDs(gomock.Any(), []int{}).Return([]domain.VideoDetail{}, nil).AnyTimes()
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:  "Success - With results",
			page:  1,
			limit: 10,
			setupMocks: func() {
				mockGrpc.EXPECT().GetByIDs(gomock.Any(), []int{1, 2}).Return([]domain.VideoDetail{
					{
						ID:          1,
						Title:       "Video title",
						Description: "Video description",
						Thumbnail:   "https://example",
						Duration:    100,
						URL:         "https://example.com",
					},
					{
						ID:          2,
						Title:       "Video title",
						Description: "Video description",
						Thumbnail:   "https://example",
						Duration:    100,
						URL:         "https://example.com",
					},
				}, nil).AnyTimes()
				mockCache.EXPECT().GetTopRankedVideos(gomock.Any(), 1, 10).Return([]int{1, 2}, nil).AnyTimes()
			},
			want: []domain.VideoDetail{
				{
					ID:          1,
					Title:       "Video title",
					Description: "Video description",
					Thumbnail:   "https://example",
					Duration:    100,
					URL:         "https://example.com",
				},
				{
					ID:          2,
					Title:       "Video title",
					Description: "Video description",
					Thumbnail:   "https://example",
					Duration:    100,
					URL:         "https://example.com",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			got, err := svc.GetTopRanked(context.Background(), tt.page, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRanking() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRanking() = %v, want %v", got, tt.want)
			}
		})
	}
}
