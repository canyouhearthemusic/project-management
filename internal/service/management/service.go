package management

import (
	"github.com/canyouhearthemusic/project-management/internal/domain/project"
	"github.com/canyouhearthemusic/project-management/internal/domain/task"
	"github.com/canyouhearthemusic/project-management/internal/domain/user"
)

type Service struct {
	userRepository    user.Repository
	taskRepository    task.Repository
	projectRepository project.Repository
}

type Configuration func(s *Service) error

func New(cfg ...Configuration) *Service {
	s := &Service{}

	for _, cfg := range cfg {
		cfg(s)
	}

	return s
}

func WithUserRepository(userRepository user.Repository) Configuration {
	return func(s *Service) error {
		s.userRepository = userRepository
		return nil
	}
}

func WithTaskRepository(taskRepository task.Repository) Configuration {
	return func(s *Service) error {
		s.taskRepository = taskRepository
		return nil
	}
}

func WithProjectRepository(projectRepository project.Repository) Configuration {
	return func(s *Service) error {
		s.projectRepository = projectRepository
		return nil
	}
}
