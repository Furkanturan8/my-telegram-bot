package telegramBots

import (
	"encoding/json"
	"fmt"
	"my-telegram-bot/helpers"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"` // Rüzgarın hızıdır
		Deg   float64 `json:"deg"`   // Rüzgarın yönünü belirtir.
	} `json:"wind"`
}

func GetWeather(city string) (string, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, os.Getenv("WEATHER_API_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("hava durumu API hatası: %s", resp.Status)
	}

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return "", err
	}

	if len(weatherResponse.Weather) == 0 {
		return "", fmt.Errorf("hava durumu bilgisi bulunamadı")
	}
	description := helpers.TranslateWeatherDescription(weatherResponse.Weather[0].Description)
	temp := weatherResponse.Main.Temp
	wind := weatherResponse.Wind

	return fmt.Sprintf("(%s) için şu anki hava durumu: \n\t Hava -> %s \n\t Sıcaklık -> %.2f°C \n\t Rüzgar -> hızı: %.2f , yönü: %.2f", city, description, temp, wind.Speed, wind.Deg), nil
}
