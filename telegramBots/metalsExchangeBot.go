package telegramBots

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

type MetalPrice struct {
	Name           string
	BuyingPrice    string
	SellingPrice   string
	PercentageDiff string
}

type Metals struct {
	Gold   []MetalPrice
	Silver []MetalPrice
}

func ConnectToExcel() (Metals, error) {
	file, err := excelize.OpenFile("data/goldDatas.xlsm")
	if err != nil {
		return Metals{}, fmt.Errorf("excel dosyası açılamadı: %v", err)
	}

	rows, err := file.GetRows("Sayfa1")
	if err != nil {
		return Metals{}, fmt.Errorf("satırlar alınamadı: %v", err)
	}

	metals := Metals{}

	for i := 2; i < len(rows); i++ { // i = 2'den başla
		row := rows[i]

		if len(row) >= 4 { // Sütun sayısının yeterli olduğunu kontrol et
			name := row[2] // Sütunlar doğru indekse göre ayarlandı
			buyingPrice := row[3]
			sellingPrice := row[4]
			percentageDiff := row[5]

			metalPrice := MetalPrice{
				Name:           name,
				BuyingPrice:    buyingPrice,
				SellingPrice:   sellingPrice,
				PercentageDiff: percentageDiff,
			}

			if strings.Contains(name, "Gümüş") || strings.Contains(name, "Platin") {
				metals.Silver = append(metals.Silver, metalPrice)
			} else {
				metals.Gold = append(metals.Gold, metalPrice)
			}
		}
	}

	return metals, nil
}

func GetMetalsValues() (string, error) {
	metals, err := ConnectToExcel()
	if err != nil {
		return "", fmt.Errorf("Veri alınamadı: %v", err)
	}

	goldFilter := []string{"Gram Altın", "Çeyrek Altın", "Yarım Altın", "Cumhuriyet Altını", "18 Ayar Bilezik", "22 Ayar Bilezik"}
	silverFilter := []string{"Gram Gümüş"}

	altin := PrintGoldPrices(metals, goldFilter)
	gumus := PrintSilverPrices(metals, silverFilter)

	return altin + gumus, nil
}

func PrintGoldPrices(metals Metals, filter []string) string {
	var result string

	result += "<b>Altın Fiyatları:</b>\n-------------------------\n"
	for _, gold := range metals.Gold {
		if contains(filter, gold.Name) {
			result += fmt.Sprintf("<b>%s:</b>\n\t <i>Alış Fiyatı:</i> %s,\n\t <i>Satış Fiyatı:</i> %s,\n\t <i>Yüzde Değişim:</i> %s\n\n",
				gold.Name, gold.BuyingPrice, gold.SellingPrice, gold.PercentageDiff)
		}
	}

	return result
}

func PrintSilverPrices(metals Metals, filter []string) string {
	var result string

	result += "<b>Gümüş Fiyatları:</b>\n-------------------------\n"
	for _, silver := range metals.Silver {
		if contains(filter, silver.Name) {
			result += fmt.Sprintf("<b>%s:</b>\n\t <i>Alış Fiyatı:</i> %s,\n\t <i>Satış Fiyatı:</i> %s,\n\t <i>Yüzde Değişim:</i> %s\n",
				silver.Name, silver.BuyingPrice, silver.SellingPrice, silver.PercentageDiff)
		}
	}

	return result
}

// filtreleme yapan yardımcı fonksiyon
func contains(filter []string, name string) bool {
	for _, f := range filter {
		if strings.Contains(name, f) {
			return true
		}
	}
	return false
}
