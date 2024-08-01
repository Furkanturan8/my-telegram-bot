package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"my-telegram-bot/models"
	"net/http"
)

type PrayerTimeService struct {
	apiBaseURL string
}

func NewPrayerTimeService(apiBaseURL string) *PrayerTimeService {
	return &PrayerTimeService{apiBaseURL: apiBaseURL}
}

// GetPrayerTimesByCity API'den namaz vakitlerini alır ve modelimize dönüştürür
func (s *PrayerTimeService) GetPrayerTimesByCity(city models.City) ([]models.PrayerTimes, error) {
	url := fmt.Sprintf("%s%s", s.apiBaseURL, city.City)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("hata 1 : ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("hata 2 : ", err)
			return nil, err
		}
		return nil, fmt.Errorf("failed to fetch data from API. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return nil, errors.New("failed to fetch data from API")
	}

	var apiResponse struct {
		Data []struct {
			Timings struct {
				Sunrise string `json:"Sunrise"`
				Dhuhr   string `json:"Dhuhr"`
				Asr     string `json:"Asr"`
				Maghrib string `json:"Maghrib"`
				Isha    string `json:"Isha"`
				Fajr    string `json:"Fajr"`
			} `json:"timings"`
			Date struct {
				Gregorian struct {
					Date  string `json:"date"`
					Day   string `json:"day"`
					Month struct {
						Number int    `json:"number"`
						En     string `json:"en"`
					} `json:"month"`
					Year string `json:"year"`
				} `json:"gregorian"`
				Hijri struct {
					Date  string `json:"date"`
					Day   string `json:"day"`
					Month struct {
						Number int    `json:"number"`
						En     string `json:"en"`
						Ar     string `json:"ar"`
					} `json:"month"`
					Year string `json:"year"`
				} `json:"hijri"`
			} `json:"date"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		fmt.Println("hata 3 : ", err)
		return nil, err
	}

	var prayerTimesList []models.PrayerTimes
	i := 1
	for _, item := range apiResponse.Data {
		prayerTimesList = append(prayerTimesList, models.PrayerTimes{
			ID:   city.ID,
			City: city.City,
			Timings: models.Timings{
				ID:              i,
				CityID:          city.ID,
				GregorianDateID: i,
				HijriDateID:     i,
				Sunrise:         item.Timings.Sunrise,
				Dhuhr:           item.Timings.Dhuhr,
				Asr:             item.Timings.Asr,
				Maghrib:         item.Timings.Maghrib,
				Isha:            item.Timings.Isha,
				Imsak:           item.Timings.Fajr,
			},
			GregorianDate: models.GregorianDate{
				Date:      item.Date.Gregorian.Date,
				Day:       item.Date.Gregorian.Day,
				Month:     item.Date.Gregorian.Month.Number,
				MonthName: item.Date.Gregorian.Month.En,
				Year:      item.Date.Gregorian.Year,
			},
			HijriDate: models.HijriDate{
				Date:  item.Date.Hijri.Date,
				Day:   item.Date.Hijri.Day,
				Month: item.Date.Hijri.Month.Number,
				Year:  item.Date.Hijri.Year,
			},
		})
		i++
	}

	return prayerTimesList, nil
}

func (s *PrayerTimeService) GetPrayerTimeByCity(city models.City, day int) (*models.PrayerTimes, error) {
	prayerTimesList, err := s.GetPrayerTimesByCity(city)
	if err != nil {
		return nil, err
	}

	// Gün 1 ile prayerTimesList uzunluğu arasında olmalıdır
	if day < 1 || day > len(prayerTimesList) {
		fmt.Printf("Invalid Day: %d\n", day)
		return nil, err
	}

	prayerTimeByDay := prayerTimesList[day-1]

	return &prayerTimeByDay, nil
}
