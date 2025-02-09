package repositories

import (
	"context"
	"fmt"

	"github.com/nghiatrann0502/trinity/internal/video/core/domain"
	"github.com/nghiatrann0502/trinity/internal/video/core/ports"
)

type memoryRepo struct {
	storage []domain.Video
}

var _ ports.Repository = (*memoryRepo)(nil)

func NewMemoryRepo() ports.Repository {
	return &memoryRepo{
		storage: []domain.Video{
			{
				ID:          1,
				Title:       "Video 1",
				Description: "Description 1",
				Thumbnail:   "Thumbnail 1",
				Duration:    100,
				URL:         "URL 1",
			}, {
				ID:          2,
				Title:       "Video 2",
				Description: "Description 2",
				Thumbnail:   "Thumbnail 2",
				Duration:    200,
				URL:         "URL 2",
			}, {
				ID:          3,
				Title:       "Video 3",
				Description: "Description 3",
				Thumbnail:   "Thumbnail 3",
				Duration:    300,
				URL:         "URL 3",
			}, {
				ID:          4,
				Title:       "Video 4",
				Description: "Description 4",
				Thumbnail:   "Thumbnail 4",
				Duration:    400,
				URL:         "URL 4",
			}, {
				ID:          5,
				Title:       "Video 5",
				Description: "Description 5",
				Thumbnail:   "Thumbnail 5",
				Duration:    500,
				URL:         "URL 5",
			},
		},
	}
}

func (m *memoryRepo) GetByID(ctx context.Context, id int) (*domain.Video, error) {
	fmt.Println("memoryRepo GetByID")
	for _, v := range m.storage {
		if v.ID == id {
			fmt.Println("fave")
			return &v, nil
		}
	}

	return nil, nil
}

func (m *memoryRepo) GetByIDs(ctx context.Context, ids []int) ([]domain.Video, error) {
	mapVideos := make(map[int]domain.Video)

	for _, v := range m.storage {
		mapVideos[v.ID] = v
	}

	videos := make([]domain.Video, len(ids))
	for i, id := range ids {
		videos[i] = mapVideos[id]
	}

	return videos, nil
}
