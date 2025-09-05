package handler

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/api/service"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/articles"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/authors"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/infra"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

type AppHandler struct {
	svc *service.Service
}

func NewAppHandler(svc *service.Service) *AppHandler {
	return &AppHandler{svc: svc}
}

// @Summary Create author
// @Tags Author
// @Accept json
// @Produce json
// @Param author body authors.AuthorInput true "Author input"
// @Success 200 {object} string "UUID"
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /author/create [post]
func (h *AppHandler) CreateAuthor(ctx context.Context, c *app.RequestContext) {
	var author authors.AuthorInput
	if err := c.Bind(&author); err != nil {
		infra.JSONError(c, 400, "Bad Request", err)
		return
	}

	id, err := h.svc.CreateAuthor(ctx, author)
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, id, "Author created successfully")
}

// @Summary Create article
// @Tags Article
// @Accept json
// @Produce json
// @Param article body articles.ArticleInput true "Article input"
// @Success 200 {object} string "UUID"
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /article/create [post]
func (h *AppHandler) CreateArticle(ctx context.Context, c *app.RequestContext) {
	var article articles.ArticleInput
	if err := c.Bind(&article); err != nil {
		infra.JSONError(c, 400, "Bad Request", err)
		return
	}

	id, err := h.svc.CreateArticle(ctx, &article)
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, id, "Article created successfully")
}

// @Summary Create many article
// @Tags Article
// @Accept json
// @Produce json
// @Param article body []articles.ArticleInput true "Article input list"
// @Success 200 {array} string "UUID list"
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /article/create-bulk [post]
func (h *AppHandler) CreateManyArticle(ctx context.Context, c *app.RequestContext) {
	var articleList []*articles.ArticleInput
	if err := c.Bind(&articleList); err != nil {
		infra.JSONError(c, 400, "Bad Request", err)
		return
	}

	idList, err := h.svc.CreateManyArticle(ctx, articleList)
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, idList, "Article created successfully")
}

// @Summary Update article
// @Tags Article
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param article body articles.ArticleInput true "Article input"
// @Success 200 {object} string "UUID"
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /article/update/{id} [put]
func (h *AppHandler) UpdateArticle(ctx context.Context, c *app.RequestContext) {
	var article articles.ArticleInput
	if err := c.Bind(&article); err != nil {
		infra.JSONError(c, 400, "Bad Request", err)
		return
	}

	idArticle := c.Param("id")
	if idArticle == "" {
		infra.JSONError(c, 400, "Missing ID", nil)
		return
	}

	id, err := h.svc.UpdateArticle(ctx, &article, uuid.MustParse(idArticle))
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, id, "Article updated successfully")
}

// @Summary Get article by key word
// @Tags Article
// @Accept json
// @Produce json
// @Param keyword query string true "Keyword"
// @Success 200 {array} articles.Article
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /article/search [get]
func (h *AppHandler) GetArticleByKeyWord(ctx context.Context, c *app.RequestContext) {
	keyword := c.Query("keyword")
	if keyword == "" {
		infra.JSONError(c, 400, "Bad Request", nil)
		return
	}

	articleList, err := h.svc.GetArticleByKeyWord(ctx, keyword)
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, articleList, "Article list")
}

// @Summary Get article with author by ID
// @Tags Article
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} articles.ArticleWithAuthor
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /article/author/{id} [get]
func (h *AppHandler) GetArticleWithAuthorByID(ctx context.Context, c *app.RequestContext) {
	idAuthor := c.Param("id")
	if idAuthor == "" {
		infra.JSONError(c, 400, "Missing Author ID", nil)
		return
	}

	articleList, err := h.svc.GetArticleWithAuthorByID(ctx, uuid.MustParse(idAuthor))
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, articleList, "Article list")
}

// @Summary Update author
// @Tags Author
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Param author body authors.AuthorInput true "Author input"
// @Success 200 {object} string "UUID"
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /author/update/{id} [put]
func (h *AppHandler) UpdateAuthor(ctx context.Context, c *app.RequestContext) {
	var author authors.AuthorInput
	if err := c.Bind(&author); err != nil {
		infra.JSONError(c, 400, "Bad Request", err)
		return
	}

	idAuthor := c.Param("id")
	if idAuthor == "" {
		infra.JSONError(c, 400, "missing author ID", nil)
		return
	}

	id, err := h.svc.UpdateAuthor(ctx, author, uuid.MustParse(idAuthor))
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, id, "Author updated successfully")
}

// @Summary Get article by author name
// @Tags Article
// @Accept json
// @Produce json
// @Param name query string true "Author name"
// @Success 200 {array} articles.ArticleWithAuthor
// @Failure 400 {object} infra.ErrorResponse
// @Failure 500 {object} infra.ErrorResponse
// @Router /article/author-name [get]
func (h *AppHandler) GetArticleByAuthorName(ctx context.Context, c *app.RequestContext) {
	name := c.Query("name")
	if name == "" {
		infra.JSONError(c, 400, "missing author name", nil)
		return
	}

	articleList, err := h.svc.GetArticleByAuthorName(ctx, name)
	if err != nil {
		infra.JSONError(c, 500, "Internal Server Error", err)
		return
	}

	infra.JSONSuccess(c, articleList, "Article list")
}
