package articles

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
)

type ArticleIndexer interface {
	Index(ctx context.Context, a *Article) error
	Search(ctx context.Context, keyword string) ([]*Article, error)
	GetAllArticle(ctx context.Context) ([]*Article, error)
	GetArticleByAuthorID(ctx context.Context, authorID uuid.UUID) ([]*Article, error)
	GetArticleByAuthorIDList(ctx context.Context, authorIDList []uuid.UUID) ([]*Article, error)
	UpdateField(ctx context.Context, id string, fields map[string]interface{}) error
}

type articleIndexer struct {
	es *elastic.Client
}

func NewArticleIndexer(es *elastic.Client) ArticleIndexer {
	return &articleIndexer{es: es}
}

func (i *articleIndexer) Index(ctx context.Context, a *Article) error {
	_, err := i.es.Index().
		Index("articles").
		Id(a.ID.String()).
		BodyJson(a).
		Do(ctx)
	return err
}

func (i *articleIndexer) GetAllArticle(ctx context.Context) ([]*Article, error) {
	searchResult, err := i.es.Search().
		Index("articles").
		Do(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]*Article, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var a Article
		if err := json.Unmarshal(hit.Source, &a); err != nil {
			continue
		}
		results = append(results, &a)
	}
	return results, nil
}

func (i *articleIndexer) GetArticleByAuthorID(ctx context.Context, authorID uuid.UUID) ([]*Article, error) {
	query := elastic.NewTermQuery("author_id.keyword", authorID.String())
	searchResult, err := i.es.Search().
		Index("articles").
		Query(query).
		Sort("created_at", false).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]*Article, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var a Article
		if err := json.Unmarshal(hit.Source, &a); err != nil {
			continue
		}
		results = append(results, &a)
	}
	return results, nil
}

func (i *articleIndexer) Search(ctx context.Context, keyword string) ([]*Article, error) {

	query := i.buildArticleWildcardQuery(keyword)

	searchResult, err := i.es.Search().
		Index("articles").
		Query(query).
		Sort("created_at", false).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*Article, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var a Article
		if err := json.Unmarshal(hit.Source, &a); err != nil {
			continue
		}
		results = append(results, &a)
	}

	return results, nil
}

func (i *articleIndexer) UpdateField(ctx context.Context, id string, fields map[string]interface{}) error {
	_, err := i.es.Update().
		Index("articles").
		Id(id).
		Doc(fields).
		DocAsUpsert(true).
		Do(ctx)
	return err
}

func (i *articleIndexer) GetArticleByAuthorIDList(ctx context.Context, authorIDList []uuid.UUID) ([]*Article, error) {
	query := elastic.NewTermsQuery("author_id.keyword", i.ChangeUIDtoInterface(authorIDList)...)

	searchResult, err := i.es.Search().
		Index("articles").
		Query(query).
		Sort("created_at", false).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*Article, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var a Article
		if err := json.Unmarshal(hit.Source, &a); err != nil {
			continue
		}
		results = append(results, &a)
	}
	return results, nil
}

func (i *articleIndexer) ChangeUIDtoInterface(arr []uuid.UUID) []interface{} {
	res := make([]interface{}, len(arr))
	for i, v := range arr {
		res[i] = v.String()
	}
	return res
}

func (i *articleIndexer) buildArticleWildcardQuery(keyword string) *elastic.BoolQuery {
	likeKeyword := "*" + keyword + "*"

	query := elastic.NewBoolQuery().
		Should(
			elastic.NewWildcardQuery("title", likeKeyword),
			elastic.NewWildcardQuery("body", likeKeyword),
		)
	return query
}
