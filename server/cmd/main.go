package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

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
	Echo         *echo.Echo
	DB           *sqlx.DB
	MinioHandler *db.MinioHandler
	GistHandler  *handlers.GistHandler
	AuthHandler  *handlers.AuthHandler
}

func NewApp(
	db *sqlx.DB,
	minioHandler *db.MinioHandler,
	GistHandler *handlers.GistHandler,
	authHandler *handlers.AuthHandler,
) *App {
	e := echo.New()
	middlewares.SetupMiddlewares(e)

	api := e.Group("/api")
	routes.SetupAPIRoutes(api, GistHandler, authHandler)

	e.Static("/public", "public")

	pages := e.Group("")
	routes.SetupPageRoutes(pages)

	if Debug := utils.GetEnv("DEBUG"); Debug == "true" {
		e.Debug = true
	}

	return &App{
		Echo:         e,
		DB:           db,
		MinioHandler: minioHandler,
		GistHandler:  GistHandler,
		AuthHandler:  authHandler,
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
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
			db.NewRetryManager,
			db.CreatePostgresConnection,
			db.CreateMinioClient,
			db.NewGistDB,
			db.NewReplyDB,
			db.NewUserDB,
			handlers.NewGistHandler,
			handlers.NewAuthHandler,
			NewApp,
		),
		fx.Invoke(RegisterHooks),
	).Run()
}
