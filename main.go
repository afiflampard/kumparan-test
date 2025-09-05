package main

import (
	"context"

	"github.com/afif-musyayyidin/hertz-boilerplate/api/router"
	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	_ "github.com/afif-musyayyidin/hertz-boilerplate/docs"
	"github.com/afif-musyayyidin/hertz-boilerplate/domain/infra"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzSwagger "github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title           Kumparan API
// @version         1.0
// @description     This is the API documentation for Kumparan Project
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.LoadConfig()
	db := infra.InitPostgres(cfg)
	ctx := context.Background()
	es := infra.ConnectElasticsearch(cfg)

	h := server.Default(server.WithHostPorts(":8080"))
	router.SetupRouter(ctx, h, db, es)
	h.GET("/swagger/*any", hertzSwagger.WrapHandler(swaggerFiles.Handler))
	h.Spin()
}
