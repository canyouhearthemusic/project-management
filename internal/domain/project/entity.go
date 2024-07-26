package project

import "github.com/canyouhearthemusic/project-management/internal/domain"

type Entity struct {
	ID          string
	Title       string
	Description string
	StartedAt   domain.OnlyDate `db:"started_at"`
	FinishedAt  domain.OnlyDate `db:"finished_at"`
	ManagerID   string          `db:"manager_id"`
}

var (
	ErrExists     = &ProjectError{"project already exists"}
	ErrNotFound   = &ProjectError{"project not found"}
	ErrSearch     = &ProjectError{"project search error"}
	ErrBadRequest = &ProjectError{"project bad request"}
)

func IsValidFilter(filter string) bool {
	if filter != "" && filter != "title" && filter != "manager" {
		return false
	}

	if filter == "" {
		return false
	}

	return true
}

type ProjectError struct {
	message string
}

func (e *ProjectError) Error() string {
	return e.message
}

func (e *ProjectError) Is(err error) bool {
	return e == err
}
