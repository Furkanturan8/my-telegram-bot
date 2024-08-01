package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"my-telegram-bot/models"
	"my-telegram-bot/services"
	"strconv"
)

type PrayerTimeHandler struct {
	Service *services.PrayerTimeService
}

func NewPrayerTimeHandler(service *services.PrayerTimeService) *PrayerTimeHandler {
	return &PrayerTimeHandler{Service: service}
}

func (h *PrayerTimeHandler) GetPrayerTimesByCity(c *fiber.Ctx) error {
	cityName := c.Params("city")

	cityID, exists := models.CityCodes[cityName]
	if !exists {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid city name")
	}

	city := models.City{
		ID:   cityID,
		City: cityName,
	}

	prayerTimes, err := h.Service.GetPrayerTimesByCity(city)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}
	return c.JSON(prayerTimes)
}

func (h *PrayerTimeHandler) GetPrayerTimeByCity(c *fiber.Ctx) error {
	cityName := c.Params("city")

	dayNumberParam := c.Params("dayNumber")
	dayNumber, err := strconv.Atoi(dayNumberParam)

	if err != nil {
		fmt.Println("dayNumber Params error : ", err)
		return err
	}

	cityID, exists := models.CityCodes[cityName]
	if !exists {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid city name")
	}

	city := models.City{
		ID:   cityID,
		City: cityName,
	}

	prayerTime, err := h.Service.GetPrayerTimeByCity(city, dayNumber)
	if err != nil || prayerTime == nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}
	return c.JSON(prayerTime)
}
