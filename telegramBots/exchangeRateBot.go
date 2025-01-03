package telegramBots

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExchangeRateResp struct {
	ExchangeRates []map[string]interface{} `json:"TCMB_AnlikKurBilgileri"`
}

func GetExchangeRate() (string, error) {
	url := fmt.Sprintf("https://hasanadiguzel.com.tr/api/kurgetir")
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("API çağrılırken hata oluştu: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API yanıtı beklenmeyen bir durum içeriyor: HTTP %d", resp.StatusCode)
	}

	var exchangeRateResp ExchangeRateResp
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRateResp); err != nil {
		return "", fmt.Errorf("döviz kurları bilgisi bulunamadı")
	}

	if len(exchangeRateResp.ExchangeRates) == 0 {
		return "", fmt.Errorf("ABD Doları ve Euro verileri bulunamadı")
	}

	// ABD Doları ve Euro için filtreleme
	var selectedRates []map[string]interface{}
	for _, rate := range exchangeRateResp.ExchangeRates {
		if rate["Isim"] == "ABD DOLARI" || rate["Isim"] == "EURO" {
			selectedRates = append(selectedRates, rate)
		}
	}

	if len(selectedRates) == 0 {
		return "", fmt.Errorf("hata5: ABD Doları ve Euro verileri bulunamadı")
	}

	// Verileri string olarak birleştirme
	result := "<b>Euro ve Dolar Kurları:</b>\n\n"
	for _, rate := range selectedRates {
		result += fmt.Sprintf("<b>%s (%s):</b>\n\t<i>Döviz Alış:</i> <b>%.4f</b>\n\t<i>Döviz Satış:</i> <b>%.4f</b>\n\t<i>Efektif Alış:</i> <b>%.4f</b>\n\t<i>Efektif Satış:</i> <b>%.4f</b>\n\n",
			rate["Isim"], rate["CurrencyName"], rate["ForexBuying"], rate["ForexSelling"], rate["BanknoteBuying"], rate["BanknoteSelling"])
	}

	return result, nil
}
