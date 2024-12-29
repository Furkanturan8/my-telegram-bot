package telegramBots

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"my-telegram-bot/handlers"
	"my-telegram-bot/helpers"
	"net/http"
	"strings"
	"time"
)

func StartTelegramBot(bot *tgbotapi.BotAPI, h *handlers.PrayerTimeHandler, GeminiApiKey string) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	geminiClient, err := NewGeminiClient(GeminiApiKey)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for update := range updates {
			if update.Message != nil {
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "<b>Merhaba, hoşgeldiniz!</b> \n\t <b>Hava durumu için:</b> \n\t Bursa (Default):  /weather \n\t Diğer Şehir: /weather şehir \n\n <b>Döviz kuruna bak:</b> \n\t Euro && Dolar: /exchange_rate \n\t Altın && Gümüş: /gold \n\n <b>Gemini ile Dil Öğren:</b>  \n\t İngilizce Öğren: /learn_english")
						msg.ParseMode = "HTML"
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Mesaj gönderim hatası: %v", err)
						}

					case "weather":
						city := "bursa" // Varsayılan şehir
						args := strings.TrimSpace(strings.TrimPrefix(update.Message.CommandArguments(), "weather"))

						if args != "" {
							city = helpers.ConvertTurkishToEnglish(strings.ToLower(args))
						}
						weatherInfo, err := GetWeather(city)
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hava durumu alınamadı (Hatalı şehir ismi girmiş olabilirsiniz!) : %v", err))
							bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo)
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
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

					case "gold":
						exc, err := GetMetalsValues()
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Gold kuru alınamadı: %v", err))
							bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, exc)
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						}

					case "learn_english":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "<b>İngilizce Öğrenme Seçenekleriniz:</b>\n\n"+
							"1. <b>Random Word:</b> /daily_word - Random kelime\n"+
							"2. <b>Aphorisms:</b> /daily_aphorisms - Özlü sözler\n"+
							"3. <b>Grammar Topics:</b> /grammar_topics - Dil bilgisi konuları\n"+
							"4. <b>Control Your Sentence:</b> /check_sentence [your sentence]- Cümlenizi kontrol edin\n")

						msg.ParseMode = "HTML"
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Mesaj gönderim hatası: %v", err)
						}

					case "grammar_topics":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID,
							"<b>Dil Bilgisi Konuları:</b>\n"+
								"1. /present_simple\n"+
								"2. /past_tense\n"+
								"3. /future_tense\n"+
								"4. /present_continuous\n"+
								"5. /past_continuous\n"+
								"6. /future_continuous\n"+
								"7. /present_perfect\n"+
								"8. /past_perfect\n"+
								"9. /future_perfect\n"+
								"10. /present_perfect_continuous\n"+
								"11. /past_perfect_continuous\n"+
								"12. /future_perfect_continuous",
						)
						msg.ParseMode = "HTML"
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Mesaj gönderim hatası: %v", err)
						}

					case "present_simple":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Present Simple", geminiClient)

					case "past_tense":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Past Tense", geminiClient)

					case "future_tense":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Future Tense", geminiClient)

					case "present_continuous":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Present Continuous", geminiClient)

					case "past_continuous":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Past Continuous", geminiClient)

					case "future_continuous":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Future Continuous", geminiClient)

					case "present_perfect":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Present Perfect", geminiClient)

					case "past_perfect":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Past Perfect", geminiClient)

					case "future_perfect":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Future Perfect", geminiClient)

					case "present_perfect_continuous":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Present Perfect Continuous", geminiClient)

					case "past_perfect_continuous":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Past Perfect Continuous", geminiClient)

					case "future_perfect_continuous":
						sendGrammarTopicMessage(bot, update.Message.Chat.ID, "Future Perfect Continuous", geminiClient)

					case "daily_word":
						// ChatGPT API'siyle günlük kelimeyi alıyoruz
						message, err := geminiClient.GetRandomWord()
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hata oluştu: %v", err))
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						} else {
							// Başarılıysa günlük kelimeyi gönderiyoruz
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", message))
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						}

					case "daily_aphorisms":
						message, err := geminiClient.GetAphorisms()
						if err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hata oluştu: %v", err))
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", message))
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						}

					case "check_sentence":
						args := strings.TrimSpace(strings.TrimPrefix(update.Message.CommandArguments(), "check_sentence"))
						if args == "" {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lütfen bir cümle yazınız")
							msg.ParseMode = "HTML"
							_, err := bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						} else {
							result, err := geminiClient.ControlSentence(args)
							if err != nil {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hata oluştu: %v", err))
								msg.ParseMode = "HTML"
								_, err := bot.Send(msg)
								if err != nil {
									log.Printf("Mesaj gönderim hatası: %v", err)
								}
							}
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", result))
							msg.ParseMode = "HTML"
							_, err = bot.Send(msg)
							if err != nil {
								log.Printf("Mesaj gönderim hatası: %v", err)
							}
						}
					}
				}
			}
		}
	}()
}

func KeepAlive() {
	for {
		resp, err := http.Get("http://localhost:3010/ping")
		if err != nil {
			log.Printf("KeepAlive isteği hatası: %v", err)
		} else {
			resp.Body.Close()
			log.Println("KeepAlive isteği başarılı")
		}

		// 10 dakikalık aralıklarla kontrol eder
		time.Sleep(10 * time.Minute)
	}
}

func sendGrammarTopicMessage(bot *tgbotapi.BotAPI, chatID int64, topic string, geminiClient *GeminiClient) {
	message, err := geminiClient.TeachGrammer(topic)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Hata oluştu: %v", err))
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Mesaj gönderim hatası: %v", err)
		}
		return
	}

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Mesaj gönderim hatası: %v", err)
	}
}
