package http

import (
	"encoding/json"
	"errors"

	"net/http"

	"github.com/canyouhearthemusic/project-management/internal/domain/project"
	"github.com/canyouhearthemusic/project-management/internal/service/management"
	"github.com/canyouhearthemusic/project-management/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProjectHandler struct {
	managementService *management.Service
}

func NewProjectHandler(service *management.Service) *ProjectHandler {
	return &ProjectHandler{
		managementService: service,
	}
}

func (h *ProjectHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.list)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
		r.Get("/tasks", h.listTasks)
	})

	r.Get("/search", h.search)

	return r
}

// create godoc
// @Summary Create a project
// @Tags Project endpoints
// @Accept json
// @Param body body project.Request true "Project request"
// @Success 201 {object} response.Response "Response"
// @Failure 400 {object} response.Response "Validation errors"
// @Router /projects [post]
func (h *ProjectHandler) create(w http.ResponseWriter, r *http.Request) {
	req := project.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = project.ErrBadRequest
		response.BadRequest(w, r, err, req)
		return
	}

	if errs := req.Validate(); errs != nil {
		errors := make([]error, len(errs))
		for i, err := range errs {
			errors[i] = err
		}

		response.BadRequests(w, r, errors, req)
		return
	}

	msg, data, err := h.managementService.CreateProject(r.Context(), req)
	if err != nil {
		response.BadRequest(w, r, err, data)
		return
	}

	response.Created(w, r, msg, data)
}

// get godoc
// @Summary Get a project
// @Tags Project endpoints
// @Accept json
// @Param id path string true "Project UUID"
// @Success 200 {object} project.Response
// @Failure 400 {string} string "Bad request"
// @Router /projects/{id} [post]
func (h *ProjectHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := h.managementService.GetProject(r.Context(), id)
	if err != nil {
		response.NotFound(w, r, err)
		return
	}

	response.OK(w, r, data)
}

// list godoc
// @Summary All projects
// @Tags Project endpoints
// @Success 200 {array} project.Response
// @Failure 400 {string} string "Bad request"
// @Router /projects [get]
func (h *ProjectHandler) list(w http.ResponseWriter, r *http.Request) {
	data, err := h.managementService.ListProjects(r.Context())
	if err != nil {
		response.BadRequest(w, r, err, err.Error())
		return
	}

	response.OK(w, r, data)
}

// @Summary Update a project
// @Tags Project endpoints
// @Accept json
// @Param id path string true "Project UUID"
// @Param body body project.UpdateRequest true "Project update request"
// @Success 200 {string} string "Project updated"
// @Failure 400 {object} []string "Validation errors"
// @Router /projects/{id} [put]
func (h *ProjectHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := project.UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NotFound(w, r, err)
		return
	}

	if errs := req.Validate(); errs != nil {
		errors := make([]error, len(errs))
		response.BadRequests(w, r, errors, req)
		return
	}

	err := h.managementService.UpdateProject(r.Context(), id, req)
	if err != nil {
		response.NotFound(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a project
// @Tags Project endpoints
// @Param id path string true "Project ID"
// @Success 200 {string} string "Project deleted"
// @Failure 400 {string} string "Bad request"
// @Router /projects/{id} [delete]
func (h *ProjectHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.managementService.DeleteProject(r.Context(), id)
	if err != nil {
		response.NotFound(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Search projects
// @Description Use either name or email query string
// @Tags Project endpoints
// @Param name query string false "Search by Name"
// @Param email query string false "Search by Email"
// @Success 200 {array} project.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /projects/search [get]
func (h *ProjectHandler) search(w http.ResponseWriter, r *http.Request) {
	var filter, val string

	for k, v := range r.URL.Query() {
		filter, val = k, v[0]
	}

	projects, err := h.managementService.SearchProjects(r.Context(), filter, val)
	if err != nil {
		if errors.Is(err, project.ErrSearch) {
			response.BadRequest(w, r, err, val)
			return
		}

		if errors.Is(err, project.ErrNotFound) {
			response.NotFound(w, r, err)
			return
		}

		response.InternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, projects)
}

// @Summary List project tasks
// @Tags Project endpoints
// @Param id path string true "Project ID"
// @Success 200 {array} task.Response
// @Failure 400 {string} string "Bad request"
// @Router /projects/{id}/tasks [get]
func (h *ProjectHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tasks, err := h.managementService.SearchTasks(r.Context(), "project_id", id)
	if err != nil {
		response.NotFound(w, r, err)
		return
	}

	render.JSON(w, r, tasks)
}
