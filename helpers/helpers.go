package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

// Türkçe karakterleri İngilizce karşılıklarına çeviren fonksiyon
func ConvertTurkishToEnglish(s string) string {
	replacements := map[string]string{
		"ç": "c", "Ç": "C",
		"ı": "i", "I": "I",
		"İ": "i",
		"ğ": "g", "Ğ": "G",
		"ö": "o", "Ö": "O",
		"ş": "s", "Ş": "S",
		"ü": "u", "Ü": "U",
	}
	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}
	fmt.Println(s)
	return s
}

func RequestLogger(c *fiber.Ctx) error {
	start := time.Now()

	// Proceed to the next middleware or handler
	err := c.Next()

	stop := time.Now()
	latency := stop.Sub(start)

	// Get the status code and method
	status := c.Response().StatusCode()
	method := c.Method()
	url := c.OriginalURL()

	// Get the client IP
	clientIP := c.IP()

	// Log format: time | status | latency | clientIP | method | url
	fmt.Printf("%s | %3d | %9v | %15s | %-7s | %s\n",
		start.Format("15:04:05"),
		status,
		latency,
		clientIP,
		method,
		url,
	)

	return err
}
