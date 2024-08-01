package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"my-telegram-bot/helpers"
	"my-telegram-bot/routes"
	"os"
	_ "strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"my-telegram-bot/handlers"
	"my-telegram-bot/services"
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

	app := fiber.New()
	app.Use(helpers.RequestLogger)

	year := time.Now().Year()
	month := int(time.Now().Month()) // Ayı tamsayı olarak almak için
	url := os.Getenv("API_BASE_URL") + fmt.Sprintf("%d/%02d?country=turkey&city=", year, month)

	prayerTimesService := services.NewPrayerTimeService(url)
	prayerTimesHandler := handlers.NewPrayerTimeHandler(prayerTimesService)
	routes.PrayerTimeRoutes(app, prayerTimesHandler)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	StartTelegramBot(bot, prayerTimesHandler)

	fmt.Println("\n\tBismillah -> Bot is running...")
	app.Listen(":3000")
}
