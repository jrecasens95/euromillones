package handlers

import (
	"errors"
	"strconv"

	"euromillones/internal/repositories"
	"euromillones/internal/services"
	"euromillones/internal/validators"

	"github.com/gofiber/fiber/v2"
)

type DrawHandler struct {
	service *services.DrawService
}

func NewDrawHandler(service *services.DrawService) *DrawHandler {
	return &DrawHandler{service: service}
}

func (h *DrawHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	draws, err := h.service.List(page, limit)
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(draws)
}

func (h *DrawHandler) Find(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, err)
	}
	draw, err := h.service.Find(id)
	if errors.Is(err, repositories.ErrDrawNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "sorteo no encontrado"})
	}
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(draw)
}

func (h *DrawHandler) Create(c *fiber.Ctx) error {
	var input validators.DrawInput
	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, err)
	}
	draw, err := h.service.Create(input)
	if err != nil {
		return badRequest(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(draw)
}

func (h *DrawHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, err)
	}

	var input validators.DrawInput
	if err := c.BodyParser(&input); err != nil {
		return badRequest(c, err)
	}
	draw, err := h.service.Update(id, input)
	if errors.Is(err, repositories.ErrDrawNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "sorteo no encontrado"})
	}
	if err != nil {
		return badRequest(c, err)
	}
	return c.JSON(draw)
}

func (h *DrawHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, err)
	}
	if err := h.service.Delete(id); errors.Is(err, repositories.ErrDrawNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "sorteo no encontrado"})
	} else if err != nil {
		return serverError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func parseID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	return uint(id), err
}

func badRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

func serverError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}
