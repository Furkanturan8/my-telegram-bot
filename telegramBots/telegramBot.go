package telegramBots

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"my-telegram-bot/handlers"
	"my-telegram-bot/helpers"
	"os"
	"strconv"
	"strings"
	"time"
)

func StartTelegramBot(bot *tgbotapi.BotAPI, h *handlers.PrayerTimeHandler) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for update := range updates {
			if update.Message != nil {
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "<b>Merhaba, hoşgeldiniz!</b> \n\n <b>Namaz vakitleri için şuna tıklayın:</b> \n\t Bursa (Default): /prayer_times \n\t Diğer Şehir: /prayer_times şehir \n\n <b>Hava durumu için:</b> \n\t Bursa (Default):  /weather \n\t Diğer Şehir: /weather şehir \n\n <b>Döviz kuruna bak:</b> \n\t Euro && Dolar: /exchange_rate")
						msg.ParseMode = "HTML"
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Mesaj gönderim hatası: %v", err)
						}

					case "prayer_times":
						city := "bursa" // Varsayılan şehir
						args := strings.TrimSpace(strings.TrimPrefix(update.Message.CommandArguments(), "prayer_times"))
						if args != "" {
							city = helpers.ConvertTurkishToEnglish(strings.ToLower(args))
						}
						SendPrayerTimes(bot, update.Message.Chat.ID, city, time.Now().Day(), h)

					case "weather":
						city := "bursa" // Varsayılan şehir
						args := strings.TrimSpace(strings.TrimPrefix(update.Message.CommandArguments(), "weather"))

						if args != "" {
							city = helpers.ConvertTurkishToEnglish(strings.ToLower(args))
						}
						weatherInfo, err := GetWeather(city)
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hava durumu alınamadı: %v", err))
							bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo)
							bot.Send(msg)
						}

					case "exchange_rate":
						exchangeRate, err := GetExchangeRate()
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Döviz kuru alınamadı: %v", err))
							bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, exchangeRate)
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						}

					default:
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bilinmeyen komut.")
						bot.Send(msg)
					}
				}
			}
		}
	}()

	go func() {
		chatID, err := strconv.Atoi(os.Getenv("CHAT_ID"))
		if err != nil {
			log.Fatal(err, "chatID hatalı!")
		}
		for {
			cityName := "bursa" // Burada varsayılan şehir adı kullanılıyor
			dayNumber := time.Now().Day()

			// Bildirim gönderme işlemini başlat
			NotifyBeforePrayer(bot, int64(chatID), cityName, dayNumber, h)

			// 1 dakikalık bir gecikme ekledik, sürekli döngüde her dakika kontrol eder
			time.Sleep(1 * time.Minute)
		}
	}()
}
