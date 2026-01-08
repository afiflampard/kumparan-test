package router

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/api/handler"
	"github.com/afif-musyayyidin/hertz-boilerplate/api/service"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/articles"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/afif-musyayyidin/hertz-boilerplate/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/jmoiron/sqlx"
	"github.com/olivere/elastic/v7"
)

func SetupRouter(ctx context.Context, h *server.Hertz, db *sqlx.DB, dbReplica *sqlx.DB, es *elastic.Client) {
	repoAuthors := authors.NewAuthorRepo(db, dbReplica)
	repoArticles := articles.NewArticleRepo(ctx, db)
	indexArticles := articles.NewArticleIndexer(es)

	authMiddleware := middleware.AuthMiddleware()
	svc := service.NewService(ctx, db, repoAuthors, repoArticles, indexArticles)
	handler := handler.NewAppHandler(svc)

	author := h.Group("/author")
	{
		author.POST("/login", handler.LoginAuthor)
		author.POST("/create", handler.CreateAuthor)
		author.PUT("/update/:id", authMiddleware, handler.UpdateAuthor)
	}
	article := h.Group("/article")
	{
		article.GET("/all", authMiddleware, handler.GetAllArticle)
		article.POST("/create", authMiddleware, handler.CreateArticle)
		article.POST("/create-bulk", authMiddleware, handler.CreateManyArticle)
		article.PUT("/update/:id", authMiddleware, handler.UpdateArticle)
		article.GET("/search", handler.GetArticleByKeyWord)
		article.GET("/author/:id", authMiddleware, handler.GetArticleWithAuthorByID)
		article.GET("/author-name", handler.GetArticleByAuthorName)
	}
}
