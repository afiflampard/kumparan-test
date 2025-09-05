package articles

import (
	"context"
	"errors"

	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/google/uuid"
)

type ArticleMutation interface {
	CreateArticle(ctx context.Context, u *ArticleInput) (*uuid.UUID, error)
	UpdateArticle(ctx context.Context, u *ArticleInput, id uuid.UUID) (*uuid.UUID, error)
	GetArticleByKeyWord(ctx context.Context, keyword string) ([]*Article, error)
	CreateManyArticle(ctx context.Context, u []*ArticleInput) ([]*uuid.UUID, error)
	GetArticleWithAuthorByID(ctx context.Context, id uuid.UUID) (*ArticleWithAuthor, error)
	GetArticleByAuthorName(ctx context.Context, name string) ([]*ArticleWithAuthor, error)
	GetArticleByID(ctx context.Context, id uuid.UUID) (*Article, error)

	Commit(ctx context.Context) error
	Cancel(ctx context.Context)
}

type articleMutation struct {
	repo   ArticleRepository
	index  ArticleIndexer
	author authors.AuthorMutation
}

func NewArticleMutation(repo ArticleRepository, index ArticleIndexer, author authors.AuthorMutation) ArticleMutation {
	return &articleMutation{repo: repo, index: index, author: author}
}

func (m *articleMutation) CreateArticle(ctx context.Context, u *ArticleInput) (*uuid.UUID, error) {
	id, err := m.repo.Save(ctx, u)
	if err != nil {
		return nil, err
	}

	if err := m.index.Index(ctx, &Article{ID: *id, Title: u.Title, Body: u.Body, AuthorID: u.AuthorID}); err != nil {
		return nil, err
	}

	return id, nil
}

func (m *articleMutation) GetArticleByID(ctx context.Context, id uuid.UUID) (*Article, error) {
	article, err := m.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	getAuthor, err := m.author.GetAuthorByID(ctx, article.AuthorID)
	if err != nil {
		return nil, err
	}
	article.Author = getAuthor
	return article, nil
}

func (m *articleMutation) UpdateArticle(ctx context.Context, u *ArticleInput, id uuid.UUID) (*uuid.UUID, error) {
	if u == nil {
		return nil, errors.New("user is nil")
	}
	if id == uuid.Nil {
		return nil, errors.New("missing user ID")
	}
	m.repo.Update(ctx, u, id)
	m.index.UpdateField(ctx, id.String(), map[string]interface{}{
		"title": u.Title,
		"body":  u.Body,
	})
	return &id, nil
}

func (m *articleMutation) GetArticleByAuthorName(ctx context.Context, name string) ([]*ArticleWithAuthor, error) {
	var (
		authorIDList               []uuid.UUID
		authorArticleWithAuthorMap = make(map[uuid.UUID]*ArticleWithAuthor)
	)
	getIDNameListAuthor, err := m.author.FindIDNameByName(ctx, name)
	if err != nil {
		return nil, err
	}
	for _, authorID := range getIDNameListAuthor {
		authorIDList = append(authorIDList, authorID.ID)
		authorArticleWithAuthorMap[authorID.ID] = &ArticleWithAuthor{
			Author: authors.Author{
				ID:   authorID.ID,
				Name: authorID.Name,
			},
		}
	}
	articleList, err := m.index.GetArticleByAuthorIDList(ctx, authorIDList)
	if err != nil {
		return nil, err
	}
	for _, article := range articleList {
		if _, ok := authorArticleWithAuthorMap[article.AuthorID]; !ok {
			continue
		}
		authorArticleWithAuthorMap[article.AuthorID].Article = append(authorArticleWithAuthorMap[article.AuthorID].Article, article)
	}
	var articleWithAuthorList []*ArticleWithAuthor
	for _, articleWithAuthor := range authorArticleWithAuthorMap {
		articleWithAuthorList = append(articleWithAuthorList, articleWithAuthor)
	}
	return articleWithAuthorList, nil
}

func (m *articleMutation) GetArticleByKeyWord(ctx context.Context, keyword string) ([]*Article, error) {
	var (
		authorList   = make(map[uuid.UUID]authors.Author)
		idAuthorList []uuid.UUID
	)
	articleList, err := m.index.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}
	for _, article := range articleList {
		idAuthorList = append(idAuthorList, article.AuthorID)
	}
	authorListResult, err := m.author.GetAuthorByIDList(ctx, idAuthorList)
	if err != nil {
		return nil, err
	}
	for _, author := range authorListResult {
		authorList[author.ID] = author
	}
	for _, article := range articleList {
		author := authorList[article.AuthorID]
		article.Author = &author
	}
	return articleList, nil
}

func (m *articleMutation) CreateManyArticle(ctx context.Context, u []*ArticleInput) ([]*uuid.UUID, error) {
	var articleListID []*uuid.UUID
	articleList, err := m.repo.CreateManyArticle(ctx, u)
	if err != nil {
		return nil, err
	}
	for _, article := range articleList {
		if err := m.index.Index(ctx, &Article{ID: article.ID, Title: article.Title, Body: article.Body, AuthorID: article.AuthorID}); err != nil {
			return nil, err
		}
	}
	for _, article := range articleList {
		articleListID = append(articleListID, &article.ID)
	}
	return articleListID, nil
}

func (m *articleMutation) GetArticleWithAuthorByID(ctx context.Context, id uuid.UUID) (*ArticleWithAuthor, error) {

	getAuthor, err := m.author.GetAuthorByID(ctx, id)
	if err != nil {
		return nil, err
	}

	getArticle, err := m.index.GetArticleByAuthorID(ctx, getAuthor.ID)
	if err != nil {
		return nil, err
	}

	// getArticle, err := m.repo.FindAllArticleWithAuthorByAuthorID(ctx, getAuthor.ID)
	// if err != nil {
	// 	return nil, err
	// }

	return &ArticleWithAuthor{
		Author:  *getAuthor,
		Article: getArticle,
	}, nil
}

// Cancel implements ArticleMutation.
func (m *articleMutation) Cancel(ctx context.Context) {
	m.repo.Cancel(ctx)
}

// Commit implements ArticleMutation.
func (m *articleMutation) Commit(ctx context.Context) error {
	return m.repo.Commit(ctx)
}
