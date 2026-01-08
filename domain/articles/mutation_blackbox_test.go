package articles_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/articles"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB *sqlx.DB
	es     *elastic.Client
	ctx    = context.Background()
)

func TestMain(m *testing.M) {
	// load config
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	var err error
	testDB, err = sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal("failed to connect test database:", err)
	}

	// elastic client
	es, err = elastic.NewClient(
		elastic.SetURL(cfg.ElasticURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal("failed to connect elastic:", err)
	}

	cleanDB()

	code := m.Run()
	os.Exit(code)
}

func cleanDB() {
	testDB.Exec("DELETE FROM articles")
	testDB.Exec("DELETE FROM authors")

	es.DeleteIndex("articles").Do(ctx)
	es.CreateIndex("articles").Do(ctx)
}

func newMutation() articles.ArticleMutation {
	ctx := context.Background()
	indexer := articles.NewArticleIndexer(es)
	repo := articles.NewArticleRepo(ctx, testDB)

	authorRepo := authors.NewAuthorRepo(testDB, testDB)
	authorMutation := authors.NewAuthorMutation(authorRepo, testDB)

	return articles.NewArticleMutation(repo, indexer, testDB, authorMutation)
}
func TestCreateArticle(t *testing.T) {
	cleanDB()
	mutation := newMutation()

	_, err := es.CreateIndex("articles").Do(ctx)
	if err != nil && !elastic.IsStatusCode(err, 400) {
		t.Fatal(err)
	}

	authorID := uuid.New()

	_, err = testDB.Exec(
		`INSERT INTO authors (id, name, email) VALUES ($1, $2, $3)`,
		authorID, "Fifa", "fifa@example.com",
	)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)

	input := articles.ArticleInput{
		Title:    "Test Article",
		Body:     "Test Body",
	}
	id, err := mutation.CreateArticle(ctx, &input, authorID)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	// verify in DB
	var title string
	err = testDB.Get(&title, "SELECT title FROM articles WHERE id = $1", *id)
	assert.NoError(t, err)
	assert.Equal(t, "Test Article", title)

	// verify in Elastic
	res, err := es.Get().
		Index("articles").
		Id(id.String()).
		Do(ctx)
	assert.NoError(t, err)
	assert.True(t, res.Found)
}

func TestUpdateArticle(t *testing.T) {
	mutation := newMutation()
	cleanDB()

	authorID := uuid.New()
	_, err := testDB.Exec(
		`INSERT INTO authors (id, name, email) VALUES ($1, $2, $3)`,
		authorID, "Dini", "dini@example.com",
	)
	assert.NoError(t, err)

	input := articles.ArticleInput{
		Title:    "Old Title",
		Body:     "Old Body",
	}
	id, err := mutation.CreateArticle(ctx, &input, authorID)
	assert.NoError(t, err)

	mutation2 := newMutation()
	updateInput := articles.ArticleInput{
		Title:    "New Title",
		Body:     "New Body",
	}
	_, err = mutation2.UpdateArticle(ctx, &updateInput, *id, authorID)
	assert.NoError(t, err)

	var title string
	err = testDB.Get(&title, "SELECT title FROM articles WHERE id = $1", *id)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", title)

	res, err := es.Get().
		Index("articles").
		Id(id.String()).
		Do(ctx)
	assert.NoError(t, err)
	assert.True(t, res.Found)
}

func TestGetArticleByID(t *testing.T) {
	mutation := newMutation()
	cleanDB()

	authorID := uuid.New()
	_, err := testDB.Exec(
		`INSERT INTO authors (id, name, email) VALUES ($1, $2, $3)`,
		authorID, "Rani", "rani@example.com",
	)
	assert.NoError(t, err)

	input := articles.ArticleInput{
		Title:    "Get Title",
		Body:     "Get Body",
	}
	id, err := mutation.CreateArticle(ctx, &input, authorID)
	assert.NoError(t, err)
	article, err := mutation.GetArticleByID(ctx, *id)
	assert.NoError(t, err)
	assert.Equal(t, "Get Title", article.Title)
	assert.Equal(t, "Rani", article.Author.Name)
}

func TestGetArticleWithAuthorByID(t *testing.T) {
	mutation := newMutation()
	cleanDB()

	authorID := uuid.New()
	_, err := testDB.Exec(
		`INSERT INTO authors (id, name, email) VALUES ($1, $2, $3)`,
		authorID, "Siti", "siti@example.com",
	)
	require.NoError(t, err)

	input := articles.ArticleInput{
		Title:    "Siti's Article",
		Body:     "Content from Siti",
	}
	_, err = mutation.CreateArticle(ctx, &input, authorID)
	require.NoError(t, err)

	_, _ = es.Refresh("articles").Do(ctx)

	result, err := mutation.GetArticleWithAuthorByID(ctx, authorID)
	require.NoError(t, err)

	assert.Equal(t, "Siti", result.Author.Name)
	assert.Len(t, result.Article, 1)
	assert.Equal(t, "Siti's Article", result.Article[0].Title)
}
