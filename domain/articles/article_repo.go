package articles

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/domain/infra/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticleRepo struct {
	db *sqlx.Tx
}

func NewArticleRepo(ctx context.Context, db *sqlx.DB) ArticleRepository {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil
	}
	return &ArticleRepo{db: tx}
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
func (a *ArticleRepo) Save(ctx context.Context, u *ArticleInput) (*uuid.UUID, error) {
	newArticle := CreateNewArticle(*u)
	_, err := a.db.NamedExecContext(ctx, CreateArticleQuery, newArticle)
	if err != nil {
		return nil, err
	}
	return &newArticle.ID, nil
}

// Update implements ArticleRepository.
func (a *ArticleRepo) Update(ctx context.Context, u *ArticleInput, id uuid.UUID) (*uuid.UUID, error) {
	article := u.ToArticleUpdate(id)
	_, err := a.db.NamedExecContext(ctx, UpdateArticleQuery, article)
	if err != nil {
		return nil, err
	}
	return &article.ID, nil
}

func (a *ArticleRepo) CreateManyArticle(ctx context.Context, u []*ArticleInput) ([]Article, error) {
	articles := CreateManyArticle(u)
	stmt, err := a.db.PrepareNamedContext(ctx, CreateArticleQuery)
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

// Cancel implements ArticleRepository.
func (a *ArticleRepo) Cancel(ctx context.Context) {
	if a != nil {
		_ = a.db.Rollback()
	}
}

// Commit implements ArticleRepository.
func (a *ArticleRepo) Commit(ctx context.Context) error {
	if a != nil {
		return a.db.Commit()
	}
	return nil
}
