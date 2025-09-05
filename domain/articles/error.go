package articles

import "github.com/afif-musyayyidin/hertz-boilerplate/domain/infra"

var (
	ErrNotFound     = infra.New("NOT_FOUND", "Not found")
	ErrInvalidInput = infra.New("INVALID_INPUT", "Invalid input")
	ErrInternal     = infra.New("INTERNAL_SERVER_ERROR", "Internal server error")
)
