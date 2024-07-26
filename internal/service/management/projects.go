package management

import (
	"context"

	"github.com/canyouhearthemusic/project-management/internal/domain"
	"github.com/canyouhearthemusic/project-management/internal/domain/project"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *Service) CreateProject(ctx context.Context, req project.Request) (string, project.Response, error) {
	logger := logrus.WithContext(ctx)

	data := project.Entity{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		ManagerID:   req.ManagerID,
		StartedAt:   domain.OnlyDate(req.StartedAt),
		FinishedAt:  domain.OnlyDate(req.FinishedAt),
	}

	msg, obj, err := s.projectRepository.Create(ctx, data)
	if err != nil {
		logger.Errorln("failed to create project")
		return "", project.Response{}, err
	}

	return msg, project.ParseFromEntity(obj), nil
}

func (s *Service) GetProject(ctx context.Context, id string) (project.Response, error) {
	logger := logrus.WithContext(ctx)

	data, err := s.projectRepository.Get(ctx, id)
	if err != nil {
		logger.Errorln("failed to get project")
		return project.Response{}, err
	}

	return project.ParseFromEntity(data), nil
}

func (s *Service) UpdateProject(ctx context.Context, id string, req project.UpdateRequest) error {
	logger := logrus.WithContext(ctx)

	data := project.Entity{
		Title:       req.Title,
		Description: req.Description,
		ManagerID:   req.ManagerID,
		FinishedAt:  domain.OnlyDate(req.FinishedAt),
	}

	err := s.projectRepository.Update(ctx, id, data)
	if err != nil {
		logger.Errorln("failed to update project")
		return err
	}

	return nil
}

func (s *Service) DeleteProject(ctx context.Context, id string) error {
	logger := logrus.WithContext(ctx)

	err := s.projectRepository.Delete(ctx, id)
	if err != nil {
		logger.Errorln("failed to delete project")
		return err
	}

	return nil
}

func (s *Service) ListProjects(ctx context.Context) ([]project.Response, error) {
	logger := logrus.WithContext(ctx)

	data, err := s.projectRepository.List(ctx)
	if err != nil {
		logger.Errorln("failed to list projects")
		return nil, project.ErrNotFound
	}

	return project.ParseFromEntities(data), nil
}

func (s *Service) SearchProjects(ctx context.Context, filter, value string) ([]project.Response, error) {
	logger := logrus.WithContext(ctx)

	if value == "" || !project.IsValidFilter(filter) {
		err := project.ErrSearch
		logger.Errorln("failed to search tasks")
		return nil, err
	}

	data, err := s.projectRepository.Search(ctx, filter, value)
	if err != nil {
		logger.Errorln("failed to search projects")
		return nil, err
	}

	return project.ParseFromEntities(data), nil
}
