package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/canyouhearthemusic/project-management/internal/domain/user"
	"github.com/canyouhearthemusic/project-management/internal/service/management"
	"github.com/canyouhearthemusic/project-management/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserHandler struct {
	managementService *management.Service
}

func NewUserHandler(service *management.Service) *UserHandler {
	return &UserHandler{managementService: service}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.create)

	r.Get("/search", h.search)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
		r.Get("/tasks", h.listTasks)
	})

	return r
}

// list godoc
// @Summary All users
// @Tags User endpoints
// @Success 200 {array} user.Response
// @Failure 400 {string} string "Bad request"
// @Router /users [get]
func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	data, err := h.managementService.ListUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.OK(w, r, data)
}

// create godoc
// @Summary Create a user
// @Tags User endpoints
// @Accept json
// @Param body body user.Request true "User request"
// @Success 201 {object} user.Response "Response"
// @Failure 400 {object} response.Response "Validation errors"
// @Router /users [post]
func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	req := user.Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NotFound(w, r, err)
		return
	}

	if errs := req.Validate(); errs != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs)
		return
	}

	msg, data, err := h.managementService.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.Created(w, r, msg, data)
}

// get godoc
// @Summary Get a user
// @Tags User endpoints
// @Accept json
// @Param id path string true "User UUID"
// @Success 201 {object} user.Response "Response"
// @Failure 400 {object} response.Response "Validation errors"
// @Router /users/{id} [get]
func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := h.managementService.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.OK(w, r, data)
}

// update godoc
// @Summary Update a user
// @Tags User endpoints
// @Accept json
// @Param id path string true "User UUID"
// @Param body body user.UpdateRequest true "User update request"
// @Success 200 {string} string "User updated"
// @Failure 400 {object} []string "Validation errors"
// @Router /users/{id} [put]
func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := user.UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NotFound(w, r, err)
		return
	}

	if errs := req.Validate(); errs != nil {
		errors := make([]error, len(errs))
		response.BadRequests(w, r, errors, req)
		return
	}

	err := h.managementService.UpdateUser(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			response.NotFound(w, r, err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// delete godoc
// @Summary Delete a user
// @Tags User endpoints
// @Param id path string true "User UUID"
// @Success 200 {string} string "User Deleted"
// @Failure 404 {object} response.Response "Not Found"
// @Router /users/{id} [delete]
func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.managementService.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			response.NotFound(w, r, err)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// list godoc
// @Summary All tasks of user
// @Tags User endpoints
// @Param id path string true "User UUID"
// @Success 200 {array} task.Response
// @Failure 400 {string} string "Bad request"
// @Router /users/{id}/tasks [get]
func (h *UserHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tasks, err := h.managementService.SearchTasks(r.Context(), "user_id", id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			response.NotFound(w, r, err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, tasks)
}

// search godoc
// @Summary Search users
// @Description You can find a users by name or email
// @Tags User endpoints
// @Param name query string false "Search by Name"
// @Param email query string false "Search by Email"
// @Success 200 {array} user.Response
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Router /users/search [get]
func (h *UserHandler) search(w http.ResponseWriter, r *http.Request) {
	var filter, val string

	for k, v := range r.URL.Query() {
		filter, val = k, v[0]
	}

	users, err := h.managementService.SearchUsers(r.Context(), filter, val)
	if err != nil {
		if errors.Is(err, user.ErrSearch) {
			response.BadRequest(w, r, err, val)
			return
		}

		if errors.Is(err, user.ErrNotFound) {
			response.NotFound(w, r, err)
			return
		}

		response.InternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, users)
}
