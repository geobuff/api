package helpers

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type key struct {
	text     string
	language string
}

var (
	cache      = make(map[key]string)
	cacheMutex = sync.RWMutex{}
)

func TranslateText(targetLanguage, text string) (string, error) {
	if val, ok := cache[key{text, targetLanguage}]; ok {
		return val, nil
	}

	ctx := context.Background()
	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}

	translatedText := resp[0].Text

	cacheMutex.Lock()
	cache[key{text, targetLanguage}] = translatedText
	cacheMutex.Unlock()

	return translatedText, nil
}
