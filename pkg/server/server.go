package server

import (
	"context"
	"net/http"
)

type Server struct {
	http *http.Server
}

type Configuration func(r *Server) error

func New(configs ...Configuration) (r *Server, err error) {
	r = &Server{}

	for _, cfg := range configs {
		if err = cfg(r); err != nil {
			return
		}
	}
	return
}

func (r *Server) Start() error {
	go func() {
		if err := r.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	return nil
}

func (r *Server) Stop(ctx context.Context) error {
	return r.http.Shutdown(ctx)
}

func WithHTTPServer(handler http.Handler, port string) Configuration {
	return func(r *Server) error {
		r.http = &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		}
		return nil
	}
}
