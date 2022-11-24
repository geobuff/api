package avatars

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/helpers"
	"github.com/geobuff/api/repo"
)

func GetAvatars(writer http.ResponseWriter, request *http.Request) {
	avatars, err := repo.GetAvatars()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	language := request.Header.Get("Content-Language")
	if language != "" && language != "en" {
		translatedAvatars := make([]repo.AvatarDto, len(avatars))
		for index, avatar := range avatars {
			avatarType, err := helpers.TranslateText(language, avatar.Type)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			description, err := helpers.TranslateText(language, avatar.Description)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			translatedAvatars[index] = repo.AvatarDto{
				ID:                avatar.ID,
				Type:              avatarType,
				CountryCode:       avatar.CountryCode,
				FlagUrl:           avatar.FlagUrl,
				Name:              avatar.Name,
				Description:       description,
				PrimaryImageUrl:   avatar.PrimaryImageUrl,
				SecondaryImageUrl: avatar.SecondaryImageUrl,
				GridPlacement:     avatar.GridPlacement,
			}
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(translatedAvatars)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(avatars)
}
