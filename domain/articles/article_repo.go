package articles

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/domain/infra/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticleRepo struct {
	db *sqlx.DB
}

func NewArticleRepo(ctx context.Context, db *sqlx.DB) ArticleRepository {
	return &ArticleRepo{db: db}
}

// FindByID implements ArticleRepository.
func (a *ArticleRepo) FindByID(ctx context.Context, id uuid.UUID) (*Article, error) {
	var article Article
	if err := a.db.GetContext(ctx, &article, FindArticleByIDQuery, id); err != nil {
		return nil, err
	}
	return &article, nil
}

// Save implements ArticleRepository.
func (a *ArticleRepo) Save(ctx context.Context, u *ArticleInput, tx *sqlx.Tx) (*Article, error) {
	newArticle := CreateNewArticle(*u)
	_, err := tx.NamedExecContext(ctx, CreateArticleQuery, newArticle)
	if err != nil {
		return nil, err
	}
	return &newArticle, nil
}

// Update implements ArticleRepository.
func (a *ArticleRepo) Update(ctx context.Context, u *ArticleInput, id uuid.UUID, tx *sqlx.Tx) (*uuid.UUID, error) {
	article := u.ToArticleUpdate(id)
	_, err := tx.NamedExecContext(ctx, UpdateArticleQuery, article)
	if err != nil {
		return nil, err
	}
	return &article.ID, nil
}

func (a *ArticleRepo) CreateManyArticle(ctx context.Context, u []*ArticleInput, tx *sqlx.Tx) ([]Article, error) {
	articles := CreateManyArticle(u)
	stmt, err := tx.PrepareNamedContext(ctx, CreateArticleQuery)
	if err != nil {
		return nil, err
	}
	for _, article := range articles {
		_, err := stmt.ExecContext(ctx, article)
		if err != nil {
			return nil, err
		}
	}

	return articles, nil
}

func (a *ArticleRepo) FindAllArticleByAuthorID(ctx context.Context, id uuid.UUID) ([]*Article, error) {
	var articles []*Article
	if err := a.db.SelectContext(ctx, &articles, FindAllArticleByAuthorIDQuery, id); err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *ArticleRepo) FindAllArticleWithAuthorByAuthorID(ctx context.Context, id uuid.UUID) ([]*Article, error) {
	var articles []*Article
	if err := a.db.SelectContext(ctx, &articles, FindAllArticleWithAuthorByAuthorIDQuery, id); err != nil {
		logger.Debug("error select article with author by author id", err)
		return nil, err
	}
	return articles, nil
}
