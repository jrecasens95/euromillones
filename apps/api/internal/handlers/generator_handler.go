package handlers

import (
	"euromillones/internal/services"

	"github.com/gofiber/fiber/v2"
)

type GeneratorHandler struct {
	service *services.GeneratorService
}

func NewGeneratorHandler(service *services.GeneratorService) *GeneratorHandler {
	return &GeneratorHandler{service: service}
}

func (h *GeneratorHandler) Generate(c *fiber.Ctx) error {
	var request services.GenerateRequest
	if err := c.BodyParser(&request); err != nil {
		return badRequest(c, err)
	}
	response, err := h.service.Generate(request)
	if err != nil {
		return badRequest(c, err)
	}
	return c.JSON(response)
}
