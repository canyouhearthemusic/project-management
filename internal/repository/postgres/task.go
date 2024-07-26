package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"database/sql"

	"github.com/canyouhearthemusic/project-management/internal/domain/task"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	if db == nil {
		panic("db is required")
	}

	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) Create(ctx context.Context, t task.Entity) (msg string, obj task.Entity, err error) {
	q := `
		INSERT INTO tasks (id, title, description, priority, status, author_id, project_id, created_at, done_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`

	args := []any{t.ID, t.Title, t.Description, t.Priority, t.Status, t.AuthorID, t.ProjectID, t.CreatedAt, t.DoneAt}

	_, err = r.db.ExecContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = task.ErrExists
		}
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return "", task.Entity{}, task.ErrExists
		}
		return
	}

	return "task has been created", t, err
}

func (r *TaskRepository) Update(ctx context.Context, id string, t task.Entity) (err error) {
	sets, args := r.prepareArgs(t)
	if len(sets) > 0 {
		args = append(args, id)
		q := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d RETURNING ID", strings.Join(sets, ", "), len(args))

		err = r.db.QueryRowContext(ctx, q, args...).Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = task.ErrNotFound
			}
		}
	}

	return
}

func (r *TaskRepository) Get(ctx context.Context, id string) (t task.Entity, err error) {
	t = task.Entity{}

	q := `
	SELECT * FROM tasks WHERE id = $1
	`

	if err = r.db.GetContext(ctx, &t, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = task.ErrNotFound
			return
		}
	}

	return
}

func (r *TaskRepository) Delete(ctx context.Context, id string) (err error) {
	q := `
	DELETE FROM tasks WHERE id = $1 RETURNING id
	`

	if err = r.db.QueryRowContext(ctx, q, id).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = task.ErrNotFound
			return
		}
	}

	return
}

func (r *TaskRepository) List(ctx context.Context) (tasks []task.Entity, err error) {
	q := "SELECT * FROM tasks"
	err = r.db.SelectContext(ctx, &tasks, q)
	if err != nil {
		return
	}

	return
}

func (r *TaskRepository) Search(ctx context.Context, filter, value string) (tasks []task.Entity, err error) {
	tasks = []task.Entity{}

	filter = r.prepareFilterArg(filter)

	q := fmt.Sprintf("SELECT * FROM users WHERE %s = $1", filter)

	err = r.db.SelectContext(ctx, &tasks, q, value)
	if err != nil {
		return
	}

	if len(tasks) == 0 {
		err = task.ErrNotFound
		return
	}

	return
}

func (r *TaskRepository) prepareArgs(data task.Entity) (sets []string, args []any) {
	if data.Title != "" {
		args = append(args, data.Title)
		sets = append(sets, fmt.Sprintf("Title=$%d", len(args)))
	}

	if data.Description != "" {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("description=$%d", len(args)))
	}

	if data.Priority != "" {
		args = append(args, data.Priority)
		sets = append(sets, fmt.Sprintf("priority=$%d", len(args)))
	}

	if data.Status != "" {
		args = append(args, data.Status)
		sets = append(sets, fmt.Sprintf("status=$%d", len(args)))
	}

	if data.AuthorID != "" {
		args = append(args, data.AuthorID)
		sets = append(sets, fmt.Sprintf("author_id=$%d", len(args)))
	}

	if data.ProjectID != "" {
		args = append(args, data.ProjectID)
		sets = append(sets, fmt.Sprintf("project_id=$%d", len(args)))
	}

	if data.DoneAt != "" {
		args = append(args, data.DoneAt)
		sets = append(sets, fmt.Sprintf("done_at=$%d", len(args)))
	}

	return
}

func (r *TaskRepository) prepareFilterArg(arg string) string {
	switch arg {
	case "title":
		return "title"
	case "description":
		return "description"
	case "priority":
		return "priority"
	case "status":
		return "status"
	case "assignee":
		return "author_id"
	case "project_id":
		return "project_id"
	case "created_at":
		return "created_at"
	case "done_at":
		return "done_at"
	default:
		return ""
	}
}
