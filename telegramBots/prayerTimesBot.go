package telegramBots

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"my-telegram-bot/handlers"
	"my-telegram-bot/helpers"
	"my-telegram-bot/models"
	"strings"
	"time"
)

var notificationsSent = make(map[string]bool) // Her namaz vakti için bildirim gönderildi mi kontrolü

func SendPrayerTimes(bot *tgbotapi.BotAPI, chatID int64, cityName string, dayNumber int, h *handlers.PrayerTimeHandler) {
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

	prayerTimesList := []struct{ name, time string }{
		{"Imsak", prayerTimes.Timings.Imsak},
		{"Gunes", prayerTimes.Timings.Sunrise},
		{"Ogle", prayerTimes.Timings.Dhuhr},
		{"Ikindi", prayerTimes.Timings.Asr},
		{"Aksam", prayerTimes.Timings.Maghrib},
		{"Yatsi", prayerTimes.Timings.Isha},
	}

	currentTime := time.Now().In(time.FixedZone("UTC+03", 3*60*60))
	var nextPrayerTime time.Time
	var prayerName string

	for _, prayerTime := range prayerTimesList {
		prayerTimeClean := strings.Split(prayerTime.time, " ")[0]
		prayerTimeDate := fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), prayerTimeClean)
		prayerTimeParsed, err := time.ParseInLocation("2006-01-02 15:04", prayerTimeDate, time.FixedZone("UTC+03", 3*60*60))
		if err != nil {
			log.Printf("Error parsing prayer time: %v", err)
			continue
		}

		if (nextPrayerTime.IsZero() || prayerTimeParsed.Before(nextPrayerTime)) && prayerTimeParsed.After(currentTime) {
			nextPrayerTime = prayerTimeParsed
			prayerName = prayerTime.name
		}
	}

	var timeLeft string
	if !nextPrayerTime.IsZero() {
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

func NotifyBeforePrayer(bot *tgbotapi.BotAPI, chatID int64, cityName string, dayNumber int, h *handlers.PrayerTimeHandler) {
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

		time.Sleep(1 * time.Minute)
	}
}
