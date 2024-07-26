package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"database/sql"

	"github.com/canyouhearthemusic/project-management/internal/domain/user"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	if db == nil {
		panic("db is required")
	}

	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, u user.Entity) (msg string, obj user.Entity, err error) {
	q := `
		INSERT INTO users (id, name, email, registration_date, role)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	args := []any{u.ID, u.Name, u.Email, u.RegistrationDate, u.Role}

	_, err = r.db.ExecContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = user.ErrExists
		}
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return "", user.Entity{}, user.ErrExists
		}
		return
	}

	return "user has been created", u, err
}

func (r *UserRepository) Update(ctx context.Context, id string, u user.Entity) (err error) {
	sets, args := r.prepareArgs(u)
	if len(sets) > 0 {
		args = append(args, id)
		q := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d RETURNING ID", strings.Join(sets, ", "), len(args))

		err = r.db.QueryRowContext(ctx, q, args...).Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = user.ErrNotFound
			}
		}
	}

	return
}

func (r *UserRepository) Get(ctx context.Context, id string) (u user.Entity, err error) {
	u = user.Entity{}

	q := `
	SELECT * FROM users WHERE id = $1
	`

	if err = r.db.GetContext(ctx, &u, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = user.ErrNotFound
			return
		}
	}

	return
}

func (r *UserRepository) Delete(ctx context.Context, id string) (err error) {
	q := `
	DELETE FROM users WHERE id = $1 RETURNING id
	`

	if err = r.db.QueryRowContext(ctx, q, id).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = user.ErrNotFound
			return
		}
	}

	return
}

func (r *UserRepository) List(ctx context.Context) (users []user.Entity, err error) {
	users = []user.Entity{}

	q := "SELECT * FROM users"

	err = r.db.SelectContext(ctx, &users, q)
	if err != nil {
		return
	}

	return
}

func (r *UserRepository) Search(ctx context.Context, filter, value string) (users []user.Entity, err error) {
	users = []user.Entity{}

	filter = r.prepareFilterArg(filter)

	q := fmt.Sprintf("SELECT * FROM users WHERE %s = $1", filter)

	err = r.db.SelectContext(ctx, &users, q, value)
	if err != nil {
		return
	}

	if len(users) == 0 {
		err = user.ErrNotFound
		return
	}

	return
}

func (r *UserRepository) prepareArgs(data user.Entity) (sets []string, args []any) {
	if data.Name != "" {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	if data.Email != "" {
		args = append(args, data.Email)
		sets = append(sets, fmt.Sprintf("email=$%d", len(args)))
	}

	if data.Role != "" {
		args = append(args, data.Role)
		sets = append(sets, fmt.Sprintf("role=$%d", len(args)))
	}

	return
}

func (r *UserRepository) prepareFilterArg(arg string) string {
	switch arg {
	case "name":
		return "name"
	case "email":
		return "email"
	case "role":
		return "role"
	default:
		return ""
	}
}
