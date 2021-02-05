// Package httperr implements a library for dealing with HTTP errors.
package httperr

import (
	"github.com/go-chi/render"
	"gopher-cache/internal/common/errors"
	"gopher-cache/internal/common/logs"
	"net/http"
)

// InternalError sends an ErrorResponse to the client with an internal error status.
func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Internal server error", http.StatusInternalServerError)
}

// Unauthorized sends an ErrorResponse to the client with an unauthorized error status.
func Unauthorised(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unauthorised", http.StatusUnauthorized)
}

// BadRequest sends an ErrorResponse to the client with a bad request error status.
func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Bad request", http.StatusBadRequest)
}

// RespondWithSlugError sends an ErrorResponse to the client based on the ErrorType of the given err.
// If err is a plain error type then the client will receive an ErrorResponse with an internal error status.
func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	slugError, ok := err.(errors.SlugError)
	if !ok {
		InternalError("internal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorised(slugError.Slug(), slugError, w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), slugError, w, r)
	default:
		InternalError(slugError.Slug(), slugError, w, r)
	}
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, logMSg string, status int) {
	logs.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMSg)
	resp := ErrorResponse{slug, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

// ErrorResponse represents the response to the client of the HTTP server.
type ErrorResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

// Render sets the HTTP status in the header.
func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
