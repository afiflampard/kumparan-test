package infra

import "encoding/json"

type APIError struct {
	Success   bool                   `json:"success"`
	ErrorCode string                 `json:"error_code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *APIError) WithDetails(details map[string]interface{}) *APIError {
	newErr := *e
	newErr.Details = details
	return &newErr
}

func (e *APIError) WithMessage(msg string) *APIError {
	newErr := *e
	newErr.Message = msg
	return &newErr
}

func New(code, msg string) *APIError {
	return &APIError{
		Success:   false,
		ErrorCode: code,
		Message:   msg,
	}
}
