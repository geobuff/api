package src

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
)

func (s *Server) getAvatars(writer http.ResponseWriter, request *http.Request) {
	avatars, err := repo.GetAvatars()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	language := request.Header.Get("Content-Language")
	if language != "" && language != "en" {
		translatedAvatars := make([]repo.AvatarDto, len(avatars))
		for index, avatar := range avatars {
			translatedAvatar, err := s.translateAvatar(avatar, language)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			translatedAvatars[index] = translatedAvatar
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(translatedAvatars)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(avatars)
}

func (s *Server) translateAvatar(avatar repo.AvatarDto, language string) (repo.AvatarDto, error) {
	avatarType, err := s.ts.TranslateText(language, avatar.Type)
	if err != nil {
		return repo.AvatarDto{}, err
	}

	description, err := s.ts.TranslateText(language, avatar.Description)
	if err != nil {
		return repo.AvatarDto{}, err
	}

	return repo.AvatarDto{
		ID:                avatar.ID,
		Type:              avatarType,
		CountryCode:       avatar.CountryCode,
		FlagUrl:           avatar.FlagUrl,
		Name:              avatar.Name,
		Description:       description,
		PrimaryImageUrl:   avatar.PrimaryImageUrl,
		SecondaryImageUrl: avatar.SecondaryImageUrl,
		GridPlacement:     avatar.GridPlacement,
	}, nil
}
