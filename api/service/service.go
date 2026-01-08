package service

import (
	"context"

	articles "github.com/afif-musyayyidin/hertz-boilerplate/domain/articles"
	authors "github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	repoAuthors  authors.AuthorRepository
	repoArticles articles.ArticleRepository
	index        articles.ArticleIndexer
	db           *sqlx.DB
}

func NewService(ctx context.Context, db *sqlx.DB, repoAuthors authors.AuthorRepository, repoArticles articles.ArticleRepository, index articles.ArticleIndexer) *Service {
	return &Service{
		repoAuthors:  repoAuthors,
		repoArticles: repoArticles,
		index:        index,
		db:           db,
	}
}

func (s *Service) CreateAuthor(ctx context.Context, u authors.AuthorInput) (*uuid.UUID, error) {
	mutation := authors.NewAuthorMutation(s.repoAuthors, s.db)
	id, err := mutation.CreateAuthor(ctx, &u)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *Service) UpdateAuthor(ctx context.Context, u authors.AuthorInput, id uuid.UUID) (*uuid.UUID, error) {
	mutation := authors.NewAuthorMutation(s.repoAuthors, s.db)
	idResult, err := mutation.UpdateAuthor(ctx, &u, id)
	if err != nil {
		return nil, err
	}
	return idResult, nil
}

func (s *Service) CreateArticle(ctx context.Context, u *articles.ArticleInput, authorID uuid.UUID) (*uuid.UUID, error) {
	mutation := articles.NewArticleMutation(s.repoArticles, s.index, s.db, authors.NewAuthorMutation(s.repoAuthors, s.db))
	idResult, err := mutation.CreateArticle(ctx, u, authorID)
	if err != nil {
		return nil, err
	}
	return idResult, nil
}

func (s *Service) CreateManyArticle(ctx context.Context, u []*articles.ArticleInput, authorID uuid.UUID) ([]*uuid.UUID, error) {
	mutation := articles.NewArticleMutation(s.repoArticles, s.index, s.db, authors.NewAuthorMutation(s.repoAuthors, s.db))
	idResult, err := mutation.CreateManyArticle(ctx, u, authorID)
	if err != nil {
		return nil, err
	}
	return idResult, nil
}

func (s *Service) UpdateArticle(ctx context.Context, u *articles.ArticleInput, id uuid.UUID, authorID uuid.UUID) (*uuid.UUID, error) {
	mutation := articles.NewArticleMutation(s.repoArticles, s.index, s.db, authors.NewAuthorMutation(s.repoAuthors, s.db))
	idResult, err := mutation.UpdateArticle(ctx, u, id, authorID)
	if err != nil {
		return nil, err
	}
	return idResult, nil
}

func (s *Service) GetArticleByKeyWord(ctx context.Context, keyword string) ([]*articles.Article, error) {
	mutation := articles.NewArticleMutation(s.repoArticles, s.index, s.db, authors.NewAuthorMutation(s.repoAuthors, s.db))
	articleList, err := mutation.GetArticleByKeyWord(ctx, keyword)
	if err != nil {
		return nil, err
	}
	return articleList, nil
}

func (s *Service) GetArticleWithAuthorByID(ctx context.Context, id uuid.UUID) (*articles.ArticleWithAuthor, error) {
	mutation := articles.NewArticleMutation(s.repoArticles, s.index, s.db, authors.NewAuthorMutation(s.repoAuthors, s.db))
	articleWithAuthor, err := mutation.GetArticleWithAuthorByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return articleWithAuthor, nil
}

func (s *Service) GetArticleByAuthorName(ctx context.Context, name string) ([]*articles.ArticleWithAuthor, error) {
	mutation := articles.NewArticleMutation(s.repoArticles, s.index, s.db, authors.NewAuthorMutation(s.repoAuthors, s.db))
	articleWithAuthorList, err := mutation.GetArticleByAuthorName(ctx, name)
	if err != nil {
		return nil, err
	}
	return articleWithAuthorList, nil
}

func (s *Service) LoginAuthor(ctx context.Context, email string, password string) (*string, error) {
	mutation := authors.NewAuthorMutation(s.repoAuthors, s.db)
	token, err := mutation.LoginAuthor(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *Service) GetAllArticle(ctx context.Context) ([]*articles.Article, error) {
	mutation := articles.NewArticleMutation(s.repoArticles, s.index, s.db, authors.NewAuthorMutation(s.repoAuthors, s.db))
	articleList, err := mutation.GetAllArticle(ctx)
	if err != nil {
		return nil, err
	}
	return articleList, nil
}