package management

import (
	"context"

	"github.com/canyouhearthemusic/project-management/internal/domain"
	"github.com/canyouhearthemusic/project-management/internal/domain/task"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *Service) CreateTask(ctx context.Context, req task.Request) (string, task.Response, error) {
	logger := logrus.WithContext(ctx)

	data := task.Entity{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		CreatedAt:   domain.OnlyDate(req.CreatedAt),
		DoneAt:      domain.OnlyDate(req.DoneAt),
		AuthorID:    req.AuthorID,
		ProjectID:   req.ProjectID,
	}

	msg, obj, err := s.taskRepository.Create(ctx, data)
	if err != nil {
		logger.Errorln("failed to create task")
		return "", task.Response{}, err
	}

	return msg, task.ParseFromEntity(obj), nil
}

func (s *Service) GetTask(ctx context.Context, id string) (task.Response, error) {
	logger := logrus.WithContext(ctx)

	data, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		logger.Errorln("failed to get task")
		return task.Response{}, err
	}

	return task.ParseFromEntity(data), nil
}

func (s *Service) UpdateTask(ctx context.Context, id string, req task.UpdateRequest) error {
	logger := logrus.WithContext(ctx)

	data := task.Entity{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		DoneAt:      domain.OnlyDate(req.DoneAt),
		AuthorID:    req.AuthorID,
		ProjectID:   req.AuthorID,
	}

	err := s.taskRepository.Update(ctx, id, data)
	if err != nil {
		logger.Errorln("failed to update task")
		return err
	}

	return nil
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	logger := logrus.WithContext(ctx)

	err := s.taskRepository.Delete(ctx, id)
	if err != nil {
		logger.Errorln("failed to delete task")
		return err
	}

	return nil
}

func (s *Service) ListTasks(ctx context.Context) ([]task.Response, error) {
	logger := logrus.WithContext(ctx)

	data, err := s.taskRepository.List(ctx)
	if err != nil {
		logger.Errorln("failed to get tasks")
		return nil, err
	}

	return task.ParseFromEntities(data), nil
}

func (s *Service) SearchTasks(ctx context.Context, filter, value string) ([]task.Response, error) {
	logger := logrus.WithContext(ctx)

	if value == "" || !task.IsValidFilter(filter) {
		err := task.ErrSearch
		logger.Errorln("failed to search tasks")
		return nil, err
	}

	data, err := s.taskRepository.Search(ctx, filter, value)
	if err != nil {
		logger.Errorln("failed to search tasks")
		return nil, err
	}

	return task.ParseFromEntities(data), nil
}
