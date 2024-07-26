package project

import (
	"time"

	"github.com/canyouhearthemusic/project-management/internal/domain"
)

type Request struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	FinishedAt  string `json:"finished_at"`
	StartedAt   string `json:"started_at"`
	ManagerID   string `json:"manager_id"`
}

type UpdateRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	FinishedAt  string `json:"finished_at,omitempty"`
	ManagerID   string `json:"manager_id,omitempty"`
}

func (p *Request) Validate() []domain.ErrorResponse {
	var errs []domain.ErrorResponse

	if len(p.Title) > 100 {
		errs = append(errs, domain.ErrorResponse{Message: "title must be less than 100 characters", Field: "title"})
	}

	if len(p.Description) > 200 {
		errs = append(errs, domain.ErrorResponse{Message: "description must be less than 200 characters", Field: "description"})
	}

	if _, err := time.Parse(domain.DateLayout, p.StartedAt); err != nil {
		errs = append(errs, domain.ErrorResponse{Message: "invalid started_at format", Field: "started_at"})
	}

	if _, err := time.Parse(domain.DateLayout, p.FinishedAt); err != nil {
		errs = append(errs, domain.ErrorResponse{Message: "invalid finished_at format", Field: "finished_at"})
	}

	return errs
}

func (p *UpdateRequest) Validate() []domain.ErrorResponse {
	var errs []domain.ErrorResponse

	if len(p.Title) > 100 && p.Title != "" {
		errs = append(errs, domain.ErrorResponse{Message: "title must be less than 100 characters", Field: "title"})
	}

	if len(p.Description) > 200 && p.Description != "" {
		errs = append(errs, domain.ErrorResponse{Message: "description must be less than 200 characters", Field: "description"})
	}

	if _, err := time.Parse(domain.DateLayout, p.FinishedAt); p.FinishedAt != "" && err != nil {
		errs = append(errs, domain.ErrorResponse{Message: "invalid finished_at format", Field: "finished_at"})
	}

	return errs
}

type Response struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FinishedAt  string `json:"finished_at"`
	StartedAt   string `json:"started_at"`
	ManagerID   string `json:"manager_id"`
}

func ParseFromEntity(p Entity) Response {
	return Response{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		FinishedAt:  p.FinishedAt.String(),
		StartedAt:   p.StartedAt.String(),
		ManagerID:   p.ManagerID,
	}
}

func ParseFromEntities(projects []Entity) []Response {
	var responses []Response
	for _, p := range projects {
		responses = append(responses, ParseFromEntity(p))
	}
	return responses
}
