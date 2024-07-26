package task

import "github.com/canyouhearthemusic/project-management/internal/domain"

type Entity struct {
	ID          string
	Title       string
	Description string
	Priority    string
	Status      string
	AuthorID    string          `db:"author_id"`
	ProjectID   string          `db:"project_id"`
	CreatedAt   domain.OnlyDate `db:"created_at"`
	DoneAt      domain.OnlyDate `db:"done_at"`
}

var (
	ErrExists     = &TaskError{"task already exists"}
	ErrNotFound   = &TaskError{"task not found"}
	ErrSearch     = &TaskError{"task search error"}
	ErrBadRequest = &TaskError{"task bad request"}
)

func IsValidFilter(filter string) bool {
	if filter == "" && filter != "title" && filter != "priority" && filter != "status" && filter != "author_id" && filter != "project_id" {
		return false
	}

	return true
}

type TaskError struct {
	message string
}

func (e *TaskError) Error() string {
	return e.message
}

func (e *TaskError) Is(err error) bool {
	return e == err
}
