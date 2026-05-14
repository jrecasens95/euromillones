package routes

import (
	"euromillones/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func Register(
	app fiber.Router,
	drawHandler *handlers.DrawHandler,
	statsHandler *handlers.StatsHandler,
	generatorHandler *handlers.GeneratorHandler,
) {
	api := app.Group("/api")

	api.Get("/draws", drawHandler.List)
	api.Get("/draws/:id", drawHandler.Find)
	api.Post("/draws", drawHandler.Create)
	api.Put("/draws/:id", drawHandler.Update)
	api.Delete("/draws/:id", drawHandler.Delete)

	api.Get("/stats/dashboard", statsHandler.Dashboard)
	api.Get("/stats/frequencies", statsHandler.Frequencies)
	api.Get("/stats/positions", statsHandler.Positions)
	api.Get("/stats/hot-cold", statsHandler.HotCold)
	api.Get("/stats/delays", statsHandler.Delays)
	api.Get("/stats/pairs", statsHandler.Pairs)

	api.Post("/generate", generatorHandler.Generate)
}
