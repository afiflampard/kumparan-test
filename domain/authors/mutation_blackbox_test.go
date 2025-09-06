package authors_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
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

func newMutation() authors.AuthorMutation {
	repo := authors.NewAuthorRepo(testDB, testDB)
	return authors.NewAuthorMutation(repo, testDB)
}
func TestCreateAuthor(t *testing.T) {
	mutation := newMutation()

	input := authors.AuthorInput{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	id, err := mutation.CreateAuthor(ctx, &input)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	var name string
	err = testDB.Get(&name, "SELECT name FROM authors WHERE id = $1", *id)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", name)
}

func TestUpdateAuthor(t *testing.T) {
	mutation := newMutation()

	authorID := uuid.New()

	_, err := testDB.Exec(
		`INSERT INTO authors (id, name, email) VALUES ($1, $2, $3)`,
		authorID, "Jane", "jane@example.com",
	)
	assert.NoError(t, err)

	input := authors.AuthorInput{
		Name:  "Jane Updated",
		Email: "jane_updated@example.com",
	}

	id, err := mutation.UpdateAuthor(ctx, &input, authorID)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	// verify in DB
	var name string
	err = testDB.Get(&name, "SELECT name FROM authors WHERE id = $1", *id)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Updated", name)
}

func TestGetAuthorByID(t *testing.T) {
	mutation := newMutation()

	authorID := uuid.New()

	_, err := testDB.Exec(
		`INSERT INTO authors (id, name, email) VALUES ($1, $2, $3)`,
		authorID, "Fifa", "fifa@example.com",
	)
	assert.NoError(t, err)

	author, err := mutation.GetAuthorByID(ctx, authorID)
	assert.NoError(t, err)
	assert.NotNil(t, author)
	assert.Equal(t, "Fifa", author.Name)
}

func TestFindIDNameByName(t *testing.T) {
	mutation := newMutation()

	authorID := uuid.New()

	_, err := testDB.Exec(
		`INSERT INTO authors (id, name, email) VALUES ($1, $2, $3)`,
		authorID, "Messi", "messi@example.com",
	)
	assert.NoError(t, err)

	result, err := mutation.FindIDNameByName(ctx, "Messi")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Messi", result[0].Name)
}
