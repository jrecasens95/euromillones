package handlers

import (
	"euromillones/internal/services"

	"github.com/gofiber/fiber/v2"
)

type StatsHandler struct {
	service *services.StatsService
}

func NewStatsHandler(service *services.StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (h *StatsHandler) Dashboard(c *fiber.Ctx) error {
	stats, err := h.service.Dashboard()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(stats)
}

func (h *StatsHandler) Frequencies(c *fiber.Ctx) error {
	stats, err := h.service.Frequencies()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(stats)
}

func (h *StatsHandler) Positions(c *fiber.Ctx) error {
	stats, err := h.service.Positions()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(stats)
}

func (h *StatsHandler) HotCold(c *fiber.Ctx) error {
	stats, err := h.service.HotCold(c.QueryInt("limit", 10))
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(stats)
}

func (h *StatsHandler) Delays(c *fiber.Ctx) error {
	stats, err := h.service.Delays()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(stats)
}

func (h *StatsHandler) Pairs(c *fiber.Ctx) error {
	stats, err := h.service.Pairs(c.QueryInt("limit", 20))
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(stats)
}
