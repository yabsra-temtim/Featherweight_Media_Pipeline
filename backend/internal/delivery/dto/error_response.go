package dto

// ErrorResponse is the standard error envelope returned by the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse builds an ErrorResponse from a plain message.
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}
