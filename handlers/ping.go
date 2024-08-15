package handlers

import (
	"github.com/gofiber/fiber/v2"
	"my-telegram-bot/services"
)

// PingHandler, sağlık kontrolü isteğini işleyen handler'dır.
type PingHandler struct {
	Service *services.PingService
}

func NewPingHandler(service *services.PingService) *PingHandler {
	return &PingHandler{
		Service: service,
	}
}

func (h *PingHandler) Ping(c *fiber.Ctx) error {
	response := h.Service.Ping()
	return c.SendString(response)
}
