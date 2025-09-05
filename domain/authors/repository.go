package authors

import (
	"context"

	"github.com/google/uuid"
)

type AuthorRepository interface {
	Save(ctx context.Context, u *AuthorInput) (*uuid.UUID, error)
	Update(ctx context.Context, u *AuthorInput, id uuid.UUID) (*uuid.UUID, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Author, error)
	FindByIDList(ctx context.Context, idList []uuid.UUID) ([]Author, error)
	FindIDNameByName(ctx context.Context, name string) ([]*AuthorIDName, error)

	Commit(ctx context.Context) error
	Cancel(ctx context.Context)
}
