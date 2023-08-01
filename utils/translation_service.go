package utils

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/language"
)

type ITranslationService interface {
	TranslateText(targetLanguage, text string) (string, error)
}

type TranslationService struct {
	cache *cache.Cache
}

func NewTranslationService() *TranslationService {
	return &TranslationService{
		cache: cache.New(cache.DefaultExpiration, cache.DefaultExpiration),
	}
}

func (t *TranslationService) TranslateText(targetLanguage, text string) (string, error) {
	key := fmt.Sprintf("%s-%s", targetLanguage, text)
	if val, found := t.cache.Get(key); found {
		return val.(string), nil
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
	t.cache.Set(key, translatedText, cache.DefaultExpiration)

	return translatedText, nil
}
