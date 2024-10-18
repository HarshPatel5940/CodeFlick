package main

import (
	"context"
	"log"
	"time"

	"github.com/HarshPatel5940/CodeFlick/internal/db"
	"github.com/HarshPatel5940/CodeFlick/internal/handlers"
	"github.com/HarshPatel5940/CodeFlick/internal/middlewares"
	"github.com/HarshPatel5940/CodeFlick/internal/routes"
	"github.com/HarshPatel5940/CodeFlick/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type App struct {
	Echo               *echo.Echo
	DB                 *sqlx.DB
	MinioHandler       *db.MinioHandler
	FileStorageHandler *handlers.FileStorageHandler
	AuthHandler        *handlers.AuthHandler
}

func NewApp(
	db *sqlx.DB,
	minioHandler *db.MinioHandler,
	fileStorageHandler *handlers.FileStorageHandler,
	authHandler *handlers.AuthHandler,
) *App {
	e := echo.New()
	middlewares.SetupMiddlewares(e)

	api := e.Group("/api")
	routes.SetupAPIRoutes(api, fileStorageHandler, authHandler)

	e.Static("/public", "public")

	pages := e.Group("")
	routes.SetupPageRoutes(pages)

	return &App{
		Echo:               e,
		DB:                 db,
		MinioHandler:       minioHandler,
		FileStorageHandler: fileStorageHandler,
		AuthHandler:        authHandler,
	}
}

func (a *App) Start(ctx context.Context) error {
	routes.StartTime = time.Now()
	address := utils.GetServerAddress()
	go func() {
		if err := a.Echo.Start(address); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}()
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	return a.Echo.Shutdown(ctx)
}

//   - To future harsh -> we are not extending RegisterHooks function here because its being
//     invoked separately which be proved by providers
func RegisterHooks(lc fx.Lifecycle, app *App) {
	lc.Append(fx.Hook{
		OnStart: app.Start,
		OnStop:  app.Stop,
	})
}

func main() {
	utils.LoadEnv()
	fx.New(
		fx.Provide(
			db.CreatePostgresConnection,
			db.CreateMinioClient,
			db.NewGistDB,
			db.NewReplyDB,
			db.NewUserDB,
			handlers.NewFilesHandler,
			handlers.NewAuthHandler,
			NewApp,
		),
		fx.Invoke(RegisterHooks),
	).Run()
}
