package translation

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/language"
)

func TranslateText(targetLanguage, text string) (string, error) {
	key := fmt.Sprintf("%s-%s", targetLanguage, text)
	if val, found := c.Get(key); found {
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
	c.Set(key, translatedText, cache.NoExpiration)

	return translatedText, nil
}
