package handler

import (
	netHttp "net/http"

	_ "github.com/canyouhearthemusic/project-management/docs"
	"github.com/canyouhearthemusic/project-management/internal/handler/http"
	"github.com/canyouhearthemusic/project-management/internal/service/management"
	"github.com/canyouhearthemusic/project-management/pkg/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Dependencies struct {
	ManagementService *management.Service
}

type Handler struct {
	deps Dependencies

	Mux *chi.Mux
}

type Configuration func(h *Handler) error

func New(deps Dependencies, cfgs ...Configuration) Handler {
	h := Handler{
		deps: deps,
	}

	for _, cfg := range cfgs {
		cfg(&h)
	}

	return h
}

// @title Project Management API
// @BasePath /api/v1/
// @version 1.0.0

// @Summary Health-Check
// @Tags Heartbeat
// @Success 200 {string} string
// @Router /heartbeat [get]
func WithHTTPHandler() Configuration {
	return func(h *Handler) error {
		h.Mux = router.New()

		userHandler := http.NewUserHandler(h.deps.ManagementService)
		taskHandler := http.NewTaskHandler(h.deps.ManagementService)
		projecthandler := http.NewProjectHandler(h.deps.ManagementService)

		h.Mux.Get("/swagger/*", httpSwagger.WrapHandler)

		h.Mux.Route("/api/v1", func(r chi.Router) {
			r.Get("/heartbeat", func(w netHttp.ResponseWriter, r *netHttp.Request) {
				render.Status(r, netHttp.StatusOK)
				render.PlainText(w, r, "OK")
			})
			r.Mount("/users", userHandler.Routes())
			r.Mount("/tasks", taskHandler.Routes())
			r.Mount("/projects", projecthandler.Routes())
		})

		return nil
	}
}
