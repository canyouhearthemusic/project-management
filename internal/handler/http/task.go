package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/canyouhearthemusic/project-management/internal/domain/task"
	"github.com/canyouhearthemusic/project-management/internal/service/management"
	"github.com/canyouhearthemusic/project-management/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TaskHandler struct {
	managementService *management.Service
}

func NewTaskHandler(service *management.Service) *TaskHandler {
	return &TaskHandler{
		managementService: service,
	}
}

func (h *TaskHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.create)
	r.Get("/", h.list)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	r.Get("/search", h.search)

	return r
}

// create godoc
// @Summary Create a task
// @Tags Task endpoints
// @Accept json
// @Param body body task.Request true "Task request"
// @Success 201 {object} response.Response "Response"
// @Failure 400 {object} response.Response "Validation errors"
// @Router /tasks [post]
func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	req := task.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = task.ErrBadRequest
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

	msg, data, err := h.managementService.CreateTask(r.Context(), req)
	if err != nil {
		response.BadRequest(w, r, err, data)
		return
	}

	response.Created(w, r, msg, data)
}

// get godoc
// @Summary Get a task
// @Tags Task endpoints
// @Accept json
// @Param id path string true "Task UUID"
// @Success 201 {object} task.Response "Response"
// @Failure 400 {object} response.Response "Validation errors"
// @Router /tasks/{id} [get]
func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := h.managementService.GetTask(r.Context(), id)
	if err != nil {
		response.NotFound(w, r, err)
		return
	}

	response.OK(w, r, data)
}

// list godoc
// @Summary All tasks
// @Tags Task endpoints
// @Success 200 {array} task.Response
// @Failure 400 {string} string "Bad request"
// @Router /tasks [get]
func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	data, err := h.managementService.ListTasks(r.Context())
	if err != nil {
		response.BadRequest(w, r, err, err.Error())
		return
	}

	response.OK(w, r, data)
}

// update godoc
// @Summary Update a task
// @Tags Task endpoints
// @Accept json
// @Param id path string true "Task UUID"
// @Param body body task.UpdateRequest true "Task update request"
// @Success 200 {string} string "Task updated"
// @Failure 400 {object} []string "Validation errors"
// @Router /tasks/{id} [put]
func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := task.UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	if errs := req.Validate(); errs != nil {
		errors := make([]error, len(errs))
		response.BadRequests(w, r, errors, req)
		return
	}

	err := h.managementService.UpdateTask(r.Context(), id, req)
	if err != nil {
		response.NotFound(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// delete godoc
// @Summary Delete a task
// @Tags Task endpoints
// @Param id path string true "Task UUID"
// @Success 200 {string} string "Task Deleted"
// @Failure 404 {object} response.Response "Not Found"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.managementService.DeleteTask(r.Context(), id)
	if err != nil {
		response.NotFound(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// search godoc
// @Summary Search tasks
// @Description You can find a tasks by title, priority, status, author_id, project_id
// @Tags Project endpoints
// @Param title query string false "Search by Title"
// @Param priority query string false "Search by Priority"
// @Param status query string false "Search by Status"
// @Param author_id query string false "Search by Author UUID"
// @Param project_id query string false "Search by Project UUID"
// @Success 200 {array} task.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /tasks/search [get]
func (h *TaskHandler) search(w http.ResponseWriter, r *http.Request) {
	var filter, val string

	for k, v := range r.URL.Query() {
		filter, val = k, v[0]
	}

	tasks, err := h.managementService.SearchTasks(r.Context(), filter, val)
	if err != nil {
		if errors.Is(err, task.ErrSearch) {
			response.BadRequest(w, r, err, val)
			return
		}

		if errors.Is(err, task.ErrNotFound) {
			response.NotFound(w, r, err)
			return
		}

		response.InternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, tasks)
}
