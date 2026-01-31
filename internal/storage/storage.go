package storage

type StatisticsStorage struct {
	Path string
}

func NewStatisticsStorage(path string) *StatisticsStorage {
	return &StatisticsStorage{Path: path}
}

func (s *StatisticsStorage) Save(stats any) error {
	return nil
}

func (s *StatisticsStorage) Load() (any, error) {
	return nil, nil
}
