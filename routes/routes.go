package routes

import (
	"github.com/gofiber/fiber/v2"
	"my-telegram-bot/handlers"
)

func PrayerTimeRoutes(app *fiber.App, prayerTimeHandler *handlers.PrayerTimeHandler) {
	app.Get("/prayer-times/:city", prayerTimeHandler.GetPrayerTimesByCity)
	app.Get("prayer-times/:city/:dayNumber", prayerTimeHandler.GetPrayerTimeByCity)
}
