package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"my-telegram-bot/handlers"
	"my-telegram-bot/helpers"
	"my-telegram-bot/routes"
	"my-telegram-bot/services"
	"my-telegram-bot/telegramBots"
	"os"
	_ "strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
	if GEMINI_API_KEY == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	app := fiber.New()
	app.Use(helpers.RequestLogger)

	year := time.Now().Year()
	month := int(time.Now().Month()) // Ayı tamsayı olarak almak için
	BASE_URL := os.Getenv("API_BASE_URL") + fmt.Sprintf("%d/%02d?country=turkey&city=", year, month)

	prayerTimesService := services.NewPrayerTimeService(BASE_URL)
	prayerTimesHandler := handlers.NewPrayerTimeHandler(prayerTimesService)
	routes.PrayerTimeRoutes(app, prayerTimesHandler)

	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)
	routes.PingRoutes(app, pingHandler)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	telegramBots.StartTelegramBot(bot, prayerTimesHandler, GEMINI_API_KEY)
	go telegramBots.KeepAlive()

	fmt.Println("\n\tBot is running...")
	app.Listen(":3010")
}
