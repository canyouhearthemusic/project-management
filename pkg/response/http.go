package response

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func OK(w http.ResponseWriter, r *http.Request, data any) {
	render.Status(r, http.StatusOK)

	v := Response{
		Success: true,
		Data:    data,
	}
	render.JSON(w, r, v)
}

func Created(w http.ResponseWriter, r *http.Request, msg string, data any) {
	render.Status(r, http.StatusCreated)

	v := Response{
		Success: true,
		Data:    data,
		Message: msg,
	}
	render.JSON(w, r, v)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error, data any) {
	render.Status(r, http.StatusBadRequest)

	v := Response{
		Success: false,
		Data:    data,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func BadRequests(w http.ResponseWriter, r *http.Request, errs []error, data any) {
	var errStr string

	for _, err := range errs {
		errStr += fmt.Sprintf("%s, ", err.Error())
	}

	render.Status(r, http.StatusBadRequest)

	v := Response{
		Success: false,
		Data:    data,
		Message: errStr,
	}
	render.JSON(w, r, v)
}

func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusNotFound)

	v := Response{
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)

	v := Response{
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}
