package service

import (
	"github.com/LazyBachelor/LazyPM/internal/models"
	"github.com/LazyBachelor/LazyPM/internal/storage"
	"errors"
)

type StatisticsService struct {
	storage *storage.Storage[models.Statistics]
}

func NewStatisticsService(storage *storage.Storage[models.Statistics]) (*StatisticsService, error) {
	if err := storage.Init(); err != nil {
		return nil, err
	}
	return &StatisticsService{
		storage: storage,
	}, nil
}

func (s *StatisticsService) GetStatistics() (models.Statistics, error) {
	if s.storage.Data == nil {
		return models.Statistics{}, errors.New("statistics data not initialized")
	}
	return *s.storage.Data, nil
}
