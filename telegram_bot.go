package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"my-telegram-bot/handlers"
	"my-telegram-bot/helpers"
	"my-telegram-bot/models"
	"os"
	"strconv"
	"strings"
	"time"
)

var notificationsSent = make(map[string]bool) // Globale alındı, Her namaz vakti için bildirim gönderildi mi kontrolü. eğer bunu yapmazsak vakit geldiğinde her dk da mesaj gönderir!

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
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Merhaba, hoşgeldiniz! \n\n Namaz vakitleri için şuna tıklayın: \n\t Bursa (Default): /prayer_times \n\t Diğer Şehir: /prayer_times_<şehir> \n\n Uyarı: şehir ismini ingilizce kelimelerle yazınız!")
						bot.Send(msg)

					case "prayer_times":
						cityName := "bursa"
						sendPrayerTimes(bot, update.Message.Chat.ID, cityName, time.Time{}.Day(), h)

					default:
						if strings.HasPrefix(update.Message.Command(), "prayer_times_") {
							cityParam := strings.TrimPrefix(update.Message.Command(), "prayer_times_")
							cityName := helpers.ConvertTurkishToEnglish(strings.ToLower(cityParam))
							sendPrayerTimes(bot, update.Message.Chat.ID, cityName, time.Time{}.Day(), h)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bilinmeyen komut.")
							bot.Send(msg)
						}
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
			notifyBeforePrayer(bot, int64(chatID), cityName, dayNumber, h)

			// 1 dakikalık bir gecikme ekledik, sürekli döngüde her dakika kontrol eder
			time.Sleep(1 * time.Minute)
		}
	}()
}

func sendPrayerTimes(bot *tgbotapi.BotAPI, chatID int64, cityName string, dayNumber int, h *handlers.PrayerTimeHandler) {
	cityID, exists := models.CityCodes[cityName]
	if !exists {
		msg := tgbotapi.NewMessage(chatID, "Geçersiz şehir ismi.")
		bot.Send(msg)
		return
	}

	city := models.City{
		ID:   cityID,
		City: cityName,
	}

	prayerTimes, err := h.Service.GetPrayerTimeByCity(city, dayNumber)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Hata: %v", err))
		bot.Send(msg)
		return
	}

	// for döngüsü için slice oluşturduk
	prayerTimesList := []struct{ name, time string }{
		{"Imsak", prayerTimes.Timings.Imsak},
		{"Gunes", prayerTimes.Timings.Sunrise},
		{"Ogle", prayerTimes.Timings.Dhuhr},
		{"Ikindi", prayerTimes.Timings.Asr},
		{"Aksam", prayerTimes.Timings.Maghrib},
		{"Yatsi", prayerTimes.Timings.Isha},
	}

	// Kalan süreyi hesaplayalım
	currentTime := time.Now().In(time.FixedZone("UTC+03", 3*60*60))
	var nextPrayerTime time.Time
	var prayerName string

	for _, prayerTime := range prayerTimesList {
		// Saat dilimi formatlarının aynı olmasına dikkat edelim
		prayerTimeClean := strings.Split(prayerTime.time, " ")[0] // Saat dilimi kısmını ayıklayın
		prayerTimeDate := fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), prayerTimeClean)
		prayerTimeParsed, err := time.ParseInLocation("2006-01-02 15:04", prayerTimeDate, time.FixedZone("UTC+03", 3*60*60))
		if err != nil {
			log.Printf("Error parsing prayer time: %v", err)
			continue
		}
		// fmt.Println("prayerTimeClean:", prayerTimeClean, "prayerTimeDate:", prayerTimeDate, "prayerTimeParsed:", prayerTimeParsed)

		// Bu kod, geçerli bir "sonraki" namaz vaktini bulmak için kullanıyoz bea!
		if (nextPrayerTime.IsZero() || prayerTimeParsed.Before(nextPrayerTime)) && prayerTimeParsed.After(currentTime) {
			nextPrayerTime = prayerTimeParsed
			fmt.Println("next:", nextPrayerTime)
			prayerName = prayerTime.name
			fmt.Println("prayerName:", prayerName)
		}
	}

	var timeLeft string
	if !nextPrayerTime.IsZero() {
		fmt.Println("current:", currentTime)
		duration := nextPrayerTime.Sub(currentTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60

		timeLeft = fmt.Sprintf("%d saat %d dakika", hours, minutes)
		if hours > 0 || minutes > 0 || seconds > 0 {
			timeLeft = fmt.Sprintf("%d saat %d dakika %d saniye", hours, minutes, seconds)
		} else {
			timeLeft = "Zaman tamamlandı."
		}
	} else {
		timeLeft = "Bugünkü namaz vakitleri tamamlandı."
	}

	response := fmt.Sprintf("Namaz Vakitleri (%s):\nImsak: %s\nGunes: %s\nOgle: %s\nIkindi: %s\nAksam: %s\nYatsi: %s\n\nBir sonraki namaz vakti (%s): %s\nKalan süre: %s",
		cityName, prayerTimes.Timings.Imsak, prayerTimes.Timings.Sunrise, prayerTimes.Timings.Dhuhr, prayerTimes.Timings.Asr, prayerTimes.Timings.Maghrib, prayerTimes.Timings.Isha,
		prayerName, nextPrayerTime.Format("15:04"), timeLeft)

	msg := tgbotapi.NewMessage(chatID, response)
	bot.Send(msg)
}

// Namaz vakitlerine 30 dakika kala bildirim gönder
func notifyBeforePrayer(bot *tgbotapi.BotAPI, chatID int64, cityName string, dayNumber int, h *handlers.PrayerTimeHandler) {
	for {
		cityName = helpers.ConvertTurkishToEnglish(strings.ToLower(cityName))

		cityID, exists := models.CityCodes[cityName]
		if !exists {
			log.Printf("Geçersiz şehir ismi: %s", cityName)
			return
		}

		city := models.City{
			ID:   cityID,
			City: cityName,
		}

		prayerTimes, err := h.Service.GetPrayerTimeByCity(city, dayNumber)
		if err != nil {
			log.Printf("Hata: %v", err)
			return
		}

		// for döngüsü için slice oluşturduk
		prayerTimesList := []struct{ name, time string }{
			{"Imsak", prayerTimes.Timings.Imsak},
			{"Gunes", prayerTimes.Timings.Sunrise},
			{"Ogle", prayerTimes.Timings.Dhuhr},
			{"Ikindi", prayerTimes.Timings.Asr},
			{"Aksam", prayerTimes.Timings.Maghrib},
			{"Yatsi", prayerTimes.Timings.Isha},
		}

		currentTime := time.Now().In(time.FixedZone("UTC+03", 3*60*60))

		for _, prayerTime := range prayerTimesList {
			prayerTimeClean := strings.Split(prayerTime.time, " ")[0]
			prayerTimeDate := fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), prayerTimeClean)
			prayerTimeParsed, err := time.ParseInLocation("2006-01-02 15:04", prayerTimeDate, time.FixedZone("UTC+03", 3*60*60))
			if err != nil {
				log.Printf("Error parsing prayer time: %v", err)
				continue
			}

			notificationTime := prayerTimeParsed.Add(-30 * time.Minute)
			if currentTime.After(notificationTime) && currentTime.Before(prayerTimeParsed) {
				if !notificationsSent[prayerTime.name] {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Namaz vaktine 30 dakika kaldı: %s (%s)", prayerTime.name, prayerTimeParsed.Format("15:04")))
					_, err := bot.Send(msg)
					if err != nil {
						log.Printf("Failed to send message: %v", err)
					} else {
						log.Printf("Notification sent for %s: %s", prayerTime.name, prayerTimeParsed.Format("15:04"))
						notificationsSent[prayerTime.name] = true
					}
				}
			}
		}

		// 1 dakikalık bir gecikme ekleyin, sürekli döngüde her dakika kontrol eder
		time.Sleep(1 * time.Minute)
	}
}
