package service

import (
	"context"

	articles "github.com/afif-musyayyidin/hertz-boilerplate/domain/articles"
	authors "github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	newMutationAuthors  authors.AuthorMutation
	newMutationArticles articles.ArticleMutation
}

func NewService(ctx context.Context, db *sqlx.DB, repoAuthors authors.AuthorRepository, repoArticles articles.ArticleRepository, index articles.ArticleIndexer) *Service {
	return &Service{
		newMutationAuthors:  authors.NewAuthorMutation(repoAuthors),
		newMutationArticles: articles.NewArticleMutation(repoArticles, index, authors.NewAuthorMutation(repoAuthors)),
	}
}

func (s *Service) CreateAuthor(ctx context.Context, u authors.AuthorInput) (*uuid.UUID, error) {
	id, err := s.newMutationAuthors.CreateAuthor(ctx, &u)
	if err != nil {
		s.newMutationAuthors.Cancel(ctx)
		return nil, err
	}
	s.newMutationAuthors.Commit(ctx)
	return id, nil
}

func (s *Service) UpdateAuthor(ctx context.Context, u authors.AuthorInput, id uuid.UUID) (*uuid.UUID, error) {
	idResult, err := s.newMutationAuthors.UpdateAuthor(ctx, &u, id)
	if err != nil {
		s.newMutationAuthors.Cancel(ctx)
		return nil, err
	}
	s.newMutationAuthors.Commit(ctx)
	return idResult, nil
}

func (s *Service) CreateArticle(ctx context.Context, u *articles.ArticleInput) (*uuid.UUID, error) {
	idResult, err := s.newMutationArticles.CreateArticle(ctx, u)
	if err != nil {
		s.newMutationArticles.Cancel(ctx)
		return nil, err
	}
	s.newMutationArticles.Commit(ctx)
	return idResult, nil
}

func (s *Service) CreateManyArticle(ctx context.Context, u []*articles.ArticleInput) ([]*uuid.UUID, error) {
	idResult, err := s.newMutationArticles.CreateManyArticle(ctx, u)
	if err != nil {
		s.newMutationArticles.Cancel(ctx)
		return nil, err
	}
	s.newMutationArticles.Commit(ctx)
	return idResult, nil
}

func (s *Service) UpdateArticle(ctx context.Context, u *articles.ArticleInput, id uuid.UUID) (*uuid.UUID, error) {
	idResult, err := s.newMutationArticles.UpdateArticle(ctx, u, id)
	if err != nil {
		s.newMutationArticles.Cancel(ctx)
		return nil, err
	}
	s.newMutationArticles.Commit(ctx)
	return idResult, nil
}

func (s *Service) GetArticleByKeyWord(ctx context.Context, keyword string) ([]*articles.Article, error) {
	return s.newMutationArticles.GetArticleByKeyWord(ctx, keyword)
}

func (s *Service) GetArticleWithAuthorByID(ctx context.Context, id uuid.UUID) (*articles.ArticleWithAuthor, error) {
	return s.newMutationArticles.GetArticleWithAuthorByID(ctx, id)
}

func (s *Service) GetArticleByAuthorName(ctx context.Context, name string) ([]*articles.ArticleWithAuthor, error) {
	return s.newMutationArticles.GetArticleByAuthorName(ctx, name)
}
