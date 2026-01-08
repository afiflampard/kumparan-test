package articles

import (
	"time"

	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Body      string    `db:"body" json:"body"`
	AuthorID  uuid.UUID `db:"author_id" json:"author_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Author *authors.Author `db:"author" json:"author"`
}

type ArticleInput struct {
	Title    string    `json:"title"`
	Body     string    `json:"body"`
}

type ArticleInputUpdate struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Body      string    `db:"body" json:"body"`
	AuthorID  uuid.UUID `db:"author_id" json:"author_id"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ArticleWithAuthor struct {
	authors.Author
	Article []*Article `json:"article"`
}

func (a *Article) TableName() string {
	return "articles"
}

func CreateNewArticle(input ArticleInput, authorID uuid.UUID) Article {
	return Article{
		ID:        uuid.New(),
		Title:     input.Title,
		Body:      input.Body,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func CreateManyArticle(input []*ArticleInput, authorID uuid.UUID) []Article {
	var articles []Article
	for _, article := range input {
		articles = append(articles, CreateNewArticle(*article, authorID))
	}
	return articles
}

func (a *ArticleInput) ToArticleUpdate(id uuid.UUID, authorID uuid.UUID) ArticleInputUpdate {
	return ArticleInputUpdate{
		ID:        id,
		Title:     a.Title,
		Body:      a.Body,
		AuthorID:  authorID,
		UpdatedAt: time.Now(),
	}
}
