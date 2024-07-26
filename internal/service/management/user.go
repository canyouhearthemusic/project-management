package management

import (
	"context"

	"github.com/canyouhearthemusic/project-management/internal/domain"
	"github.com/canyouhearthemusic/project-management/internal/domain/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *Service) ListUsers(ctx context.Context) ([]user.Response, error) {
	logger := logrus.WithContext(ctx)

	data, err := s.userRepository.List(ctx)
	if err != nil {
		logger.Errorln("failed to get users")
		return nil, err
	}

	return user.ParseFromEntities(data), nil
}

func (s *Service) CreateUser(ctx context.Context, req user.Request) (string, user.Response, error) {
	logger := logrus.WithContext(ctx)

	data := user.Entity{
		ID:               uuid.NewString(),
		Name:             req.Name,
		Email:            req.Email,
		RegistrationDate: domain.OnlyDate(req.RegistrationDate),
		Role:             req.Role,
	}

	msg, obj, err := s.userRepository.Create(ctx, data)
	if err != nil {
		logger.Errorln("failed to create user")
		return "", user.Response{}, err
	}

	return msg, user.ParseFromEntity(obj), nil
}

func (s *Service) GetUser(ctx context.Context, id string) (user.Response, error) {
	logger := logrus.WithContext(ctx)

	data, err := s.userRepository.Get(ctx, id)
	if err != nil {
		logger.Errorln("failed to get user")
		return user.Response{}, err
	}

	return user.ParseFromEntity(data), nil
}

func (s *Service) UpdateUser(ctx context.Context, id string, req user.UpdateRequest) error {
	logger := logrus.WithContext(ctx)

	data := user.Entity{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	err := s.userRepository.Update(ctx, id, data)
	if err != nil {
		logger.Errorln("failed to update user")
		return err
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	logger := logrus.WithContext(ctx)

	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		logger.Errorln("failed to delete user")
		return err
	}

	return nil
}

func (s *Service) SearchUsers(ctx context.Context, filter, value string) ([]user.Response, error) {
	logger := logrus.WithContext(ctx)

	if value == "" || !user.IsValidFilter(filter) {
		err := user.ErrSearch
		logger.Errorln("failed to search users")
		return nil, err
	}

	data, err := s.userRepository.Search(ctx, filter, value)
	if err != nil {
		logger.Errorln("failed to search users")
		return nil, err
	}

	return user.ParseFromEntities(data), nil
}
