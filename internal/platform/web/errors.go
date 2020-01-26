package web

// ErrorResponse is the form used for API responses form failures in the API
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to pass an error during the request through the
// application with web specific context
type Error struct {
	Err    error
	Status int
}

// NewRequestError wraps a provided error with an HTTP status code.
// This function should be used when handlers gets some error
func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

// Error implements the error inrerface. It uses the default message of the
// Wrapped error. This is the error which will be shown in services error log
func (err *Error) Error() string {
	return err.Err.Error()
}
