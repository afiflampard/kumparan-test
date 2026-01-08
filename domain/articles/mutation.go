package articles

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticleMutation interface {
	CreateArticle(ctx context.Context, u *ArticleInput, authorID uuid.UUID) (*uuid.UUID, error)
	UpdateArticle(ctx context.Context, u *ArticleInput, id uuid.UUID, authorID uuid.UUID) (*uuid.UUID, error)
	GetArticleByKeyWord(ctx context.Context, keyword string) ([]*Article, error)
	CreateManyArticle(ctx context.Context, u []*ArticleInput, authorID uuid.UUID) ([]*uuid.UUID, error)
	GetArticleWithAuthorByID(ctx context.Context, id uuid.UUID) (*ArticleWithAuthor, error)
	GetArticleByAuthorName(ctx context.Context, name string) ([]*ArticleWithAuthor, error)
	GetArticleByID(ctx context.Context, id uuid.UUID) (*Article, error)
	GetAllArticle(ctx context.Context) ([]*Article, error)
}

type articleMutation struct {
	repo   ArticleRepository
	index  ArticleIndexer
	db     *sqlx.DB
	author authors.AuthorMutation
}

func NewArticleMutation(repo ArticleRepository, index ArticleIndexer, db *sqlx.DB, author authors.AuthorMutation) ArticleMutation {
	return &articleMutation{
		repo:   repo,
		index:  index,
		db:     db,
		author: author,
	}
}

func (m *articleMutation) CreateArticle(ctx context.Context, u *ArticleInput, authorID uuid.UUID) (*uuid.UUID, error) {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	newArticle, err := m.repo.Save(ctx, u, authorID, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := m.index.Index(ctx, newArticle); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &newArticle.ID, nil
}

func (m *articleMutation) GetArticleByID(ctx context.Context, id uuid.UUID) (*Article, error) {
	article, err := m.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"id": id,
		})
	}
	getAuthor, err := m.author.GetAuthorByID(ctx, article.AuthorID)
	if err != nil {
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"id": article.AuthorID,
		})
	}
	article.Author = getAuthor
	return article, nil
}

func (m *articleMutation) UpdateArticle(ctx context.Context, u *ArticleInput, id uuid.UUID, authorID uuid.UUID) (*uuid.UUID, error) {
	if u == nil {
		return nil, ErrInvalidInput.WithDetails(map[string]interface{}{
			"id": id,
		})
	}
	if id == uuid.Nil {
		return nil, ErrInvalidInput.WithDetails(map[string]interface{}{
			"id": id,
		})
	}
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	idResult, err := m.repo.Update(ctx, u, id, authorID, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := m.index.UpdateField(ctx, id.String(), map[string]interface{}{
		"title": u.Title,
		"body":  u.Body,
	}); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return idResult, nil
}

func (m *articleMutation) GetArticleByAuthorName(ctx context.Context, name string) ([]*ArticleWithAuthor, error) {
	var (
		authorIDList               []uuid.UUID
		authorArticleWithAuthorMap = make(map[uuid.UUID]*ArticleWithAuthor)
	)
	getIDNameListAuthor, err := m.author.FindIDNameByName(ctx, name)
	if err != nil {
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"name": name,
		})
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
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"authorIDList": authorIDList,
		})
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
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"keyword": keyword,
		})
	}
	for _, article := range articleList {
		idAuthorList = append(idAuthorList, article.AuthorID)
	}
	authorListResult, err := m.author.GetAuthorByIDList(ctx, idAuthorList)
	if err != nil {
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"id_author_list": idAuthorList,
		})
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

func (m *articleMutation) CreateManyArticle(ctx context.Context, u []*ArticleInput, authorID uuid.UUID) ([]*uuid.UUID, error) {
	var articleListID []*uuid.UUID
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	articleList, err := m.repo.CreateManyArticle(ctx, u, authorID, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"article_list": articleList,
		})
	}
	for _, article := range articleList {
		if err := m.index.Index(ctx, &article); err != nil {
			_ = tx.Rollback()
			return nil, ErrNotFound.WithDetails(map[string]interface{}{
				"article": article,
			})
		}
	}
	for _, article := range articleList {
		articleListID = append(articleListID, &article.ID)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return articleListID, nil
}

func (m *articleMutation) GetArticleWithAuthorByID(ctx context.Context, id uuid.UUID) (*ArticleWithAuthor, error) {

	getAuthor, err := m.author.GetAuthorByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"id": id,
		})
	}

	getArticle, err := m.index.GetArticleByAuthorID(ctx, getAuthor.ID)
	if err != nil {
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"id_author": getAuthor.ID,
		})
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

func (m *articleMutation) GetAllArticle(ctx context.Context) ([]*Article, error) {
	articleList, err := m.index.GetAllArticle(ctx)
	if err != nil {
		return nil, ErrNotFound.WithDetails(map[string]interface{}{
			"error": err,
		})
	}
	for _, article := range articleList {
		getAuthor, err := m.author.GetAuthorByID(ctx, article.AuthorID)
		if err != nil {
			return nil, ErrNotFound.WithDetails(map[string]interface{}{
				"id": article.AuthorID,
			})
		}
		article.Author = getAuthor
	}
	return articleList, nil
}
