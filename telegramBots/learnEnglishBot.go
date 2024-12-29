package telegramBots

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"strings"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	// Google PaLM API istemcisi oluşturma
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	// Oturum açan client'ı döndür
	return &GeminiClient{client: client}, nil
}

func (g *GeminiClient) generateContent(prompt string) (string, error) {
	model := g.client.GenerativeModel("models/gemini-pro")

	// API'ye istek gönderme
	resp, err := model.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Gemini API: %v", err)
	}

	// API'den gelen yanıtı işleme
	content := string(resp.Candidates[0].Content.Parts[0].(genai.Text))
	log.Println("Received content from Gemini API:", content)

	return content, nil
}

func (g *GeminiClient) TeachGrammer(request string) (string, error) {
	prompt := fmt.Sprintf("(NOTE: Please don't use markdown syntax in your explanation for example star. Don't use table.) I am learning a new language. My main language is Turkish. Teach me both english and turkish descriptions with examples this topic: " + request)
	content, err := g.generateContent(prompt)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (g *GeminiClient) ControlSentence(sentence string) (string, error) {
	prompt := fmt.Sprintf("I am learning a new language. Please teach me if my sentence is wrong, and explain what is wrong with it. Also, provide explanations in both English and Turkish for better understanding. Correct the following sentence if necessary: " + sentence)
	content, err := g.generateContent(prompt)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (g *GeminiClient) GetRandomWord() (string, error) {
	// İstek prompt'u
	prompt := "Give a random English word of the day with its meaning and example usage in English. example usage -> Word: Go (Gitmek) | Meaning: To go somewhere (bir yere gitmek) | Sentence: I am going now (Gidiyorum şuan)"

	// API'den içeriği al
	content, err := g.generateContent(prompt)
	if err != nil {
		return "", err
	}

	result := strings.Replace(content, "|", "\n\n", -1)
	message := fmt.Sprintf("<b>  --Random Word-- </b>\n\n" + result)

	return message, nil
}

func (g *GeminiClient) GetAphorisms() (string, error) {
	// İstek prompt'u
	prompt := "Give me a random English aphorism, but add a Turkish explanation! Example usage -> Aphorism: A bird in the hand is worth two in the bush. | Turkish Explanation: Eldeki bir kuş, çalıdaki iki kuştan iyidir. (Anlamı: Kesin olan bir şey, belirsiz olan iki şeyden daha değerlidir.)"

	// API'den içeriği al
	content, err := g.generateContent(prompt)
	if err != nil {
		return "", err
	}

	result := strings.Replace(content, "|", "\n\n", -1)
	message := fmt.Sprintf("<b>  --Aphorisms-- </b>\n\n" + result)

	return message, nil
}
