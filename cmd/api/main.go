package main

import (
	"context"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/db"
	"github.com/HarshPatel5940/CodeFlick/internal/handlers"
	"github.com/HarshPatel5940/CodeFlick/internal/middlewares"
	"github.com/HarshPatel5940/CodeFlick/internal/routes"
	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"go.uber.org/fx"
)

func CreateServer(
	FileStorageHandler *handlers.FileStorageHandler, AuthHandler *handlers.AuthHandler, db *sqlx.DB, minio *minio.Client,
) *echo.Echo {

	app := echo.New()
	// TODO: add middleware for route body validator
	middlewares.SetupMiddlewares(app)

	api := app.Group("/api")
	routes.SetupRoutes(api, FileStorageHandler, AuthHandler)

	return app
}

func InitServer(lifecycle fx.Lifecycle, app *echo.Echo) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			routes.StartTime = time.Now()
			address := utils.GetServerAddress()
			go app.Start(address)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown(ctx)
		},
	})
}

func main() {
	fx.New(
		fx.Provide(
			handlers.NewFilesHandler,
			handlers.NewAuthHandler,
			db.CreatePostgresConnection,
			db.CreateMinioClient,
			CreateServer,
		),
		fx.Invoke(
			utils.LoadEnv,
			db.InitMinioClient,
			InitServer,
		),
	).Run()
}
