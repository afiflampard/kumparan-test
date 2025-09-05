package authors

import "github.com/afif-musyayyidin/hertz-boilerplate/domain/infra"

var (
	ErrInvalidInput = infra.New("INVALID_INPUT", "Invalid input")
	ErrNotFound     = infra.New("NOT_FOUND", "Not found")
	ErrInternal     = infra.New("INTERNAL_SERVER_ERROR", "Internal server error")
)
