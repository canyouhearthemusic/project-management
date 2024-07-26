package user

import "github.com/canyouhearthemusic/project-management/internal/domain"

type Entity struct {
	ID               string
	Name             string
	Email            string
	RegistrationDate domain.OnlyDate `db:"registration_date"`
	Role             string
}

var (
	ErrExists     = &UserError{"user already exists"}
	ErrNotFound   = &UserError{"user not found"}
	ErrSearch     = &UserError{"user search error"}
	ErrBadRequest = &UserError{"user bad request"}
)

func IsValidFilter(filter string) bool {
	if filter == "" && filter != "name" && filter != "email" {
		return false
	}

	return true
}

type UserError struct {
	message string
}

func (e *UserError) Error() string {
	return e.message
}

func (e *UserError) Is(err error) bool {
	return e == err
}
