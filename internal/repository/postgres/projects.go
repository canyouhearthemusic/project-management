package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/canyouhearthemusic/project-management/internal/domain/project"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	if db == nil {
		panic("db is required")
	}

	return &ProjectRepository{
		db: db,
	}
}

func (r *ProjectRepository) Create(ctx context.Context, p project.Entity) (string, project.Entity, error) {
	q := `
		INSERT INTO projects (id, title, description, manager_id, started_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	args := []any{p.ID, p.Title, p.Description, p.ManagerID, p.StartedAt, p.FinishedAt}

	_, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", project.Entity{}, project.ErrExists
		}
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return "", project.Entity{}, project.ErrExists
		}
		return "", project.Entity{}, err
	}

	return "project has been created", p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, id string, p project.Entity) (err error) {
	sets, args := r.prepareArgs(p)
	if len(sets) > 0 {
		args = append(args, id)
		q := fmt.Sprintf("UPDATE projects SET %s WHERE id = $%d RETURNING ID", strings.Join(sets, ", "), len(args))

		err = r.db.QueryRowContext(ctx, q, args...).Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = project.ErrNotFound
			}
		}
	}

	return
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) (err error) {
	q := `
	DELETE FROM projects
	WHERE id = $1 RETURNING id
	`

	err = r.db.QueryRowContext(ctx, q, id).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = project.ErrNotFound
		}
	}

	return
}

func (r *ProjectRepository) Get(ctx context.Context, id string) (p project.Entity, err error) {
	p = project.Entity{}

	q := `
	SELECT * FROM projects WHERE id = $1
	`

	err = r.db.GetContext(ctx, &p, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = project.ErrNotFound
		}
	}

	return
}

func (r *ProjectRepository) List(ctx context.Context) (projects []project.Entity, err error) {
	s := "SELECT * FROM projects"

	err = r.db.SelectContext(ctx, &projects, s)
	if err != nil {
		return
	}

	return
}

func (r *ProjectRepository) Search(ctx context.Context, arg, value string) (projects []project.Entity, err error) {
	projects = []project.Entity{}

	filter := r.prepareFilterArg(arg)

	q := fmt.Sprintf("SELECT * FROM projects WHERE %s = $1", filter)

	err = r.db.SelectContext(ctx, &projects, q, value)
	if err != nil {
		return
	}

	if len(projects) == 0 {
		err = project.ErrNotFound
		return
	}

	return
}

func (r *ProjectRepository) prepareArgs(p project.Entity) (sets []string, args []any) {
	if p.Title != "" {
		args = append(args, p.Title)
		sets = append(sets, fmt.Sprintf("title = $%d", len(args)))
	}

	if p.Description != "" {
		args = append(args, p.Description)
		sets = append(sets, fmt.Sprintf("description = $%d", len(args)))
	}

	if p.ManagerID != "" {
		args = append(args, p.ManagerID)
		sets = append(sets, fmt.Sprintf("manager_id = $%d", len(args)))
	}

	if p.StartedAt != "" {
		args = append(args, p.StartedAt)
		sets = append(sets, fmt.Sprintf("started_at = $%d", len(args)))
	}

	if p.FinishedAt != "" {
		args = append(args, p.FinishedAt)
		sets = append(sets, fmt.Sprintf("finished_at = $%d", len(args)))
	}

	return
}

func (r *ProjectRepository) prepareFilterArg(arg string) string {
	switch arg {
	case "title":
		return "title"
	case "manager":
		return "manager_id"
	default:
		return ""
	}
}
