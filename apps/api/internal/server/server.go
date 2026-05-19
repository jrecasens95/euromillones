package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	healthhttp "euromillones/internal/features/health/http"
	"euromillones/internal/handlers"
	"euromillones/internal/platform/config"
	"euromillones/internal/repositories"
	"euromillones/internal/routes"
	"euromillones/internal/services"
)

func New() (*fiber.App, error) {
	drawRepo, err := repositories.NewDrawRepository(config.Current.DrawsDataDir)
	if err != nil {
		return nil, err
	}

	drawService := services.NewDrawService(drawRepo)
	statsService := services.NewStatsService(drawService)
	generatorService := services.NewGeneratorService(statsService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Current.CORSAllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	healthhttp.RegisterRoutes(app)
	routes.Register(
		app,
		handlers.NewDrawHandler(drawService),
		handlers.NewStatsHandler(statsService),
		handlers.NewGeneratorHandler(generatorService),
	)

	return app, nil
}
