package main

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/api/router"
	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/infra"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	cfg := config.LoadConfig()
	db := infra.InitPostgres(cfg)
	ctx := context.Background()
	es := infra.ConnectElasticsearch(cfg)

	h := server.Default(server.WithHostPorts(":8080"))
	router.SetupRouter(ctx, h, db, es)
	h.Spin()
}
