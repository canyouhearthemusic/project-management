package repository

import (
	"github.com/canyouhearthemusic/project-management/config"
	"github.com/canyouhearthemusic/project-management/internal/domain/project"
	"github.com/canyouhearthemusic/project-management/internal/domain/task"
	"github.com/canyouhearthemusic/project-management/internal/domain/user"
	"github.com/canyouhearthemusic/project-management/internal/repository/postgres"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres postgres.DB

	User    user.Repository
	Task    task.Repository
	Project project.Repository
}

func New(configs ...Configuration) (*Repository, error) {
	repo := &Repository{}

	for _, cfg := range configs {
		if err := cfg(repo); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

func WithPostgresStore(cfg config.DB) Configuration {
	return func(repo *Repository) (err error) {
		repo.postgres, err = postgres.New(cfg)
		if err != nil {
			return
		}

		err = repo.postgres.Migrate()
		if err != nil {
			return
		}

		repo.User = postgres.NewUserRepository(repo.postgres.Client)
		repo.Task = postgres.NewTaskRepository(repo.postgres.Client)
		repo.Project = postgres.NewProjectRepository(repo.postgres.Client)

		return
	}
}
