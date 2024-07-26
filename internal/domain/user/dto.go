package user

import (
	"regexp"
	"time"

	"github.com/canyouhearthemusic/project-management/internal/domain"
)

type Request struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	RegistrationDate string `json:"registration_date"`
}

type UpdateRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
}

func (u *Request) Validate() []domain.ErrorResponse {
	var errs []domain.ErrorResponse

	if u.Name == "" {
		errs = append(errs, domain.ErrorResponse{Message: "name is required", Field: "name"})
	}

	if ok, _ := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, u.Email); !ok {
		errs = append(errs, domain.ErrorResponse{Message: "invalid email address", Field: "email"})
	}

	if u.Role != "admin" && u.Role != "manager" {
		errs = append(errs, domain.ErrorResponse{Message: "invalid role", Field: "role"})
	}

	if _, err := time.Parse(domain.DateLayout, u.RegistrationDate); err != nil {
		errs = append(errs, domain.ErrorResponse{Message: "invalid registration_date format", Field: "registration_date"})
	}

	return errs
}

func (u *UpdateRequest) Validate() []domain.ErrorResponse {
	var errs []domain.ErrorResponse

	if ok, _ := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, u.Email); u.Email != "" && !ok {
		errs = append(errs, domain.ErrorResponse{Message: "invalid email address", Field: "email"})
	}

	if u.Role != "" && u.Role != "admin" && u.Role != "manager" {
		errs = append(errs, domain.ErrorResponse{Message: "invalid role", Field: "role"})
	}

	return errs
}

type Response struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	RegistrationDate string `json:"registration_date"`
}

func ParseFromEntity(u Entity) Response {
	return Response{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		Role:             u.Role,
		RegistrationDate: u.RegistrationDate.String(),
	}
}

func ParseFromEntities(users []Entity) []Response {
	var res []Response

	for _, u := range users {
		res = append(res, ParseFromEntity(u))
	}

	return res
}
