package user

import "context"

type Storage interface {
	CreateUser(ctx context.Context) error
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

type Service struct {
	storage Storage
}

func (s *Service) Auth(ctx context.Context) error {
	return nil
}
