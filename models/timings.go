package models

// Timings modeli
type Timings struct {
	ID              int    `json:"id"`
	CityID          int    `json:"city_id"`
	GregorianDateID int    `json:"gregorian_date_id"` // miladi takvim
	HijriDateID     int    `json:"hijri_date_id"`     // hicri takvim
	Imsak           string `json:"imsak"`             // imsak
	Sunrise         string `json:"sunrise"`           // gün doğumu
	Dhuhr           string `json:"dhuhr"`             // öğle
	Asr             string `json:"asr"`               // ikindi
	Maghrib         string `json:"maghrib"`           // aksam
	Isha            string `json:"isha"`              // yatsı
}

// GregorianDate modeli
type GregorianDate struct {
	Date      string `json:"date"`
	Day       string `json:"day"`
	Month     int    `json:"month"`
	MonthName string `json:"month_name"`
	Year      string `json:"year"`
}

// HijriDate modeli
type HijriDate struct {
	Date  string `json:"date"`
	Day   string `json:"day"`
	Month int    `json:"month"`
	Year  string `json:"year"`
}

// PrayerTimes modeli, hem namaz saatlerini hem de tarihleri içerir --> Yani yukarıdakilerin hepsini kapsar
type PrayerTimes struct {
	ID            int           `json:"id"`
	City          string        `json:"city"`
	Timings       Timings       `json:"timings"`
	GregorianDate GregorianDate `json:"gregorian_date"`
	HijriDate     HijriDate     `json:"hijri_date"`
}
