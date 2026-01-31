package service

import "beadstest/internal/storage"

type StatisticsService struct {
	storage *storage.StatisticsStorage
}

func NewStatisticsService(storage *storage.StatisticsStorage) *StatisticsService {
	return &StatisticsService{
		storage: storage,
	}
}
