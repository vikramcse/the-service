package web

// ErrorResponse is the form used for API responses form failures in the API
type ErrorResponse struct {
	Error string `json:"error"`
}
