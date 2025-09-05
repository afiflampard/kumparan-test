package articles

import (
	"context"

	"github.com/google/uuid"
)

type ArticleRepository interface {
	Save(ctx context.Context, u *ArticleInput) (*uuid.UUID, error)
	Update(ctx context.Context, u *ArticleInput, id uuid.UUID) (*uuid.UUID, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Article, error)
	FindAllArticleByAuthorID(ctx context.Context, id uuid.UUID) ([]*Article, error)
	FindAllArticleWithAuthorByAuthorID(ctx context.Context, id uuid.UUID) ([]*Article, error)
	CreateManyArticle(ctx context.Context, u []*ArticleInput) ([]Article, error)
	Commit(ctx context.Context) error
	Cancel(ctx context.Context)
}
