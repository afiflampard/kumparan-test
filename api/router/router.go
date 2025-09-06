package router

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/api/handler"
	"github.com/afif-musyayyidin/hertz-boilerplate/api/service"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/articles"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/jmoiron/sqlx"
	"github.com/olivere/elastic/v7"
)

func SetupRouter(ctx context.Context, h *server.Hertz, db *sqlx.DB, dbReplica *sqlx.DB, es *elastic.Client) {
	repoAuthors := authors.NewAuthorRepo(db, dbReplica)
	repoArticles := articles.NewArticleRepo(ctx, db)
	indexArticles := articles.NewArticleIndexer(es)

	svc := service.NewService(ctx, db, repoAuthors, repoArticles, indexArticles)
	handler := handler.NewAppHandler(svc)

	author := h.Group("/author")
	{
		author.POST("/create", handler.CreateAuthor)
		author.PUT("/update/:id", handler.UpdateAuthor)
	}
	article := h.Group("/article")
	{
		article.POST("/create", handler.CreateArticle)
		article.POST("/create-bulk", handler.CreateManyArticle)
		article.PUT("/update/:id", handler.UpdateArticle)
		article.GET("/search", handler.GetArticleByKeyWord)
		article.GET("/author/:id", handler.GetArticleWithAuthorByID)
		article.GET("/author-name", handler.GetArticleByAuthorName)
	}
}
