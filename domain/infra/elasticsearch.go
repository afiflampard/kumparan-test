package infra

import (
	"log"

	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	"github.com/olivere/elastic/v7"
)

func ConnectElasticsearch(cfg config.Config) *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(cfg.ElasticURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}
	return client
}
