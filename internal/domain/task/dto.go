package task

import (
	"time"

	"github.com/canyouhearthemusic/project-management/internal/domain"
)

type Request struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	AuthorID    string `json:"author_id"`
	ProjectID   string `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	DoneAt      string `json:"done_at"`
}

type UpdateRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Status      string `json:"status,omitempty"`
	AuthorID    string `json:"author_id,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	DoneAt      string `json:"done_at,omitempty"`
}

func (t *Request) Validate() []domain.ErrorResponse {
	var errs []domain.ErrorResponse

	if t.Title == "" {
		errs = append(errs, domain.ErrorResponse{Message: "title is required", Field: "title"})
	}
	if len(t.Title) > 100 {
		errs = append(errs, domain.ErrorResponse{Message: "title must be less than 100 characters", Field: "title"})
	}

	if t.Description == "" {
		errs = append(errs, domain.ErrorResponse{Message: "description is required", Field: "description"})
	}
	if len(t.Description) >= 200 {
		errs = append(errs, domain.ErrorResponse{Message: "description must be less than 200 characters", Field: "description"})
	}

	if _, err := time.Parse(domain.DateLayout, t.CreatedAt); err != nil {
		errs = append(errs, domain.ErrorResponse{Message: "invalid created_at format", Field: "created_at"})
	}

	if _, err := time.Parse(domain.DateLayout, t.DoneAt); err != nil {
		errs = append(errs, domain.ErrorResponse{Message: "invalid done_at format", Field: "done_at"})
	}

	if !isValidPriority(t.Priority) {
		errs = append(errs, domain.ErrorResponse{Message: "invalid priority value", Field: "priority"})
	}

	if !isValidStatus(t.Status) {
		errs = append(errs, domain.ErrorResponse{Message: "invalid status value", Field: "status"})
	}

	return errs
}

func isValidPriority(priority string) bool {
	allowedPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}
	return allowedPriorities[priority]
}

func isValidStatus(status string) bool {
	allowedStatuses := map[string]bool{
		"active":      true,
		"in_progress": true,
		"done":        true,
	}
	return allowedStatuses[status]
}

func (t *UpdateRequest) Validate() []domain.ErrorResponse {
	var errs []domain.ErrorResponse

	if len(t.Title) > 100 && t.Title != "" {
		errs = append(errs, domain.ErrorResponse{Message: "title must be less than 100 characters", Field: "title"})
	}

	if len(t.Description) >= 200 && t.Description != "" {
		errs = append(errs, domain.ErrorResponse{Message: "description must be less than 200 characters", Field: "description"})
	}

	if _, err := time.Parse(domain.DateLayout, t.DoneAt); t.DoneAt != "" && err != nil {
		errs = append(errs, domain.ErrorResponse{Message: "invalid done_at format", Field: "done_at"})
	}

	if t.Priority != "" && !isValidPriority(t.Priority) {
		errs = append(errs, domain.ErrorResponse{Message: "invalid priority value", Field: "priority"})
	}

	if t.Status != "" && !isValidStatus(t.Status) {
		errs = append(errs, domain.ErrorResponse{Message: "invalid status value", Field: "status"})
	}

	return errs
}

type Response struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	AuthorID    string `json:"author_id"`
	ProjectID   string `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	DoneAt      string `json:"done_at"`
}

func ParseFromEntity(t Entity) Response {
	return Response{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    t.Priority,
		Status:      t.Status,
		AuthorID:    t.AuthorID,
		ProjectID:   t.ProjectID,
		CreatedAt:   t.CreatedAt.String(),
		DoneAt:      t.DoneAt.String(),
	}
}

func ParseFromEntities(tasks []Entity) []Response {
	var responses []Response
	for _, t := range tasks {
		responses = append(responses, ParseFromEntity(t))
	}
	return responses
}
