package articles

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticleRepository interface {
	Save(ctx context.Context, u *ArticleInput, authorID uuid.UUID, tx *sqlx.Tx) (*Article, error)
	Update(ctx context.Context, u *ArticleInput, id uuid.UUID, authorID uuid.UUID, tx *sqlx.Tx) (*uuid.UUID, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Article, error)
	FindAllArticleByAuthorID(ctx context.Context, id uuid.UUID) ([]*Article, error)
	FindAllArticleWithAuthorByAuthorID(ctx context.Context, id uuid.UUID) ([]*Article, error)
	CreateManyArticle(ctx context.Context, u []*ArticleInput, authorID uuid.UUID, tx *sqlx.Tx) ([]Article, error)
}
