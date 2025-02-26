// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ranking/core/ports/repository.go
//
// Generated by this command:
//
//	mockgen --source=internal/ranking/core/ports/repository.go --destination=internal/ranking/adapters/repositories/mysql.go.mock.go -package=repositories
//

// Package repositories is a generated GoMock package.
package repositories

import (
	context "context"
	reflect "reflect"

	domain "github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateVideoInteraction mocks base method.
func (m *MockRepository) CreateVideoInteraction(ctx context.Context, dto domain.VideoRankingUpdate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVideoInteraction", ctx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateVideoInteraction indicates an expected call of CreateVideoInteraction.
func (mr *MockRepositoryMockRecorder) CreateVideoInteraction(ctx, dto any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVideoInteraction", reflect.TypeOf((*MockRepository)(nil).CreateVideoInteraction), ctx, dto)
}

// GetVideoRankingByVideoId mocks base method.
func (m *MockRepository) GetVideoRankingByVideoId(ctx context.Context, videoId int) (*domain.VideoRanking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideoRankingByVideoId", ctx, videoId)
	ret0, _ := ret[0].(*domain.VideoRanking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideoRankingByVideoId indicates an expected call of GetVideoRankingByVideoId.
func (mr *MockRepositoryMockRecorder) GetVideoRankingByVideoId(ctx, videoId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideoRankingByVideoId", reflect.TypeOf((*MockRepository)(nil).GetVideoRankingByVideoId), ctx, videoId)
}
