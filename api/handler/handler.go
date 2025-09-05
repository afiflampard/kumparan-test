package handler

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/api/service"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/articles"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

type AppHandler struct {
	svc *service.Service
}

func NewAppHandler(svc *service.Service) *AppHandler {
	return &AppHandler{svc: svc}
}

func (h *AppHandler) CreateAuthor(ctx context.Context, c *app.RequestContext) {
	var author authors.AuthorInput
	if err := c.Bind(&author); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateAuthor(ctx, author)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(201, map[string]interface{}{"id": id})
}

func (h *AppHandler) CreateArticle(ctx context.Context, c *app.RequestContext) {
	var article articles.ArticleInput
	if err := c.Bind(&article); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateArticle(ctx, &article)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(201, map[string]interface{}{"id": id})
}

func (h *AppHandler) CreateManyArticle(ctx context.Context, c *app.RequestContext) {
	var articleList []*articles.ArticleInput
	if err := c.Bind(&articleList); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	idList, err := h.svc.CreateManyArticle(ctx, articleList)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(201, map[string]interface{}{"id_list": idList})
}

func (h *AppHandler) UpdateArticle(ctx context.Context, c *app.RequestContext) {
	var article articles.ArticleInput
	if err := c.Bind(&article); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	idArticle := c.Param("id")
	if idArticle == "" {
		c.JSON(400, map[string]string{"error": "missing article ID"})
		return
	}

	id, err := h.svc.UpdateArticle(ctx, &article, uuid.MustParse(idArticle))
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(201, map[string]interface{}{"id": id})
}

func (h *AppHandler) GetArticleByKeyWord(ctx context.Context, c *app.RequestContext) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(400, map[string]string{"error": "missing keyword"})
		return
	}

	articleList, err := h.svc.GetArticleByKeyWord(ctx, keyword)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"article_list": articleList})
}

func (h *AppHandler) GetArticleWithAuthorByID(ctx context.Context, c *app.RequestContext) {
	idAuthor := c.Param("id")
	if idAuthor == "" {
		c.JSON(400, map[string]string{"error": "missing author ID"})
		return
	}

	articleList, err := h.svc.GetArticleWithAuthorByID(ctx, uuid.MustParse(idAuthor))
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"data": articleList})
}

func (h *AppHandler) UpdateAuthor(ctx context.Context, c *app.RequestContext) {
	var author authors.AuthorInput
	if err := c.Bind(&author); err != nil {
		c.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	idAuthor := c.Param("id")
	if idAuthor == "" {
		c.JSON(400, map[string]string{"error": "missing author ID"})
		return
	}

	id, err := h.svc.UpdateAuthor(ctx, author, uuid.MustParse(idAuthor))
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(201, map[string]interface{}{"id": id})
}

func (h *AppHandler) GetArticleByAuthorName(ctx context.Context, c *app.RequestContext) {
	name := c.Query("name")
	if name == "" {
		c.JSON(400, map[string]string{"error": "missing author name"})
		return
	}

	articleList, err := h.svc.GetArticleByAuthorName(ctx, name)
	if err != nil {
		c.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"article_list": articleList})
}
