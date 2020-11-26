// Package errors implements library for creating common errors.
package errors

// ErrorType indicates the type of an error.
type ErrorType struct {
	t string
}

// ErrorTypes used by this package.
var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
)

// SlugError extends error. In addition to error it provides a slug.
type SlugError struct {
	error     string
	slug      string
	errorType ErrorType
}

// Error returns a string indicating what error occured.
func (s SlugError) Error() string {
	return s.error
}

// Slug returns a slug string indicating the type of the error.
func (s SlugError) Slug() string {
	return s.slug
}

// ErrorType returns an ErrorType to assist in determining the type of the error.
func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}

// NewSlugError creates a new SlugError.
func NewSlugError(error string, slug string) SlugError {
	return SlugError{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeUnknown,
	}
}

// NewAuthorizationError creates a new SlugError for authorization errors.
func NewAuthorizationError(error string, slug string) SlugError {
	return SlugError{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeAuthorization,
	}
}

// NewIncorrectInputError creates a new SlugError for client input errors.
func NewIncorrectInputError(error string, slug string) SlugError {
	return SlugError{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeIncorrectInput,
	}
}
