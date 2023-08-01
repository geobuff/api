package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

type QuizPageDto struct {
	Quizzes []repo.Quiz `json:"quizzes"`
	HasMore bool        `json:"hasMore"`
}

func (s *Server) getQuizzes(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var filter repo.QuizzesFilterDto
	err = json.Unmarshal(requestBody, &filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	quizzes, err := repo.GetQuizzes(filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	language := request.Header.Get("Content-Language")
	if language != "" && language != "en" {
		translatedQuizzes := make([]repo.Quiz, len(quizzes))

		for index, quiz := range quizzes {
			translatedQuiz, err := s.translateQuiz(quiz, language)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			translatedQuizzes[index] = translatedQuiz
		}

		switch _, err := repo.GetFirstQuizID((filter.Page + 1) * filter.Limit); err {
		case sql.ErrNoRows:
			entriesDto := QuizPageDto{translatedQuizzes, false}
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(entriesDto)
		case nil:
			entriesDto := QuizPageDto{translatedQuizzes, true}
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(entriesDto)
		default:
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		}
		return
	}

	switch _, err := repo.GetFirstQuizID((filter.Page + 1) * filter.Limit); err {
	case sql.ErrNoRows:
		entriesDto := QuizPageDto{quizzes, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := QuizPageDto{quizzes, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

func (s *Server) translateQuiz(quiz repo.Quiz, language string) (repo.Quiz, error) {
	name, err := s.ts.TranslateText(language, quiz.Name)
	if err != nil {
		return repo.Quiz{}, err
	}

	plural, err := s.ts.TranslateText(language, quiz.Plural)
	if err != nil {
		return repo.Quiz{}, err
	}

	return repo.Quiz{
		ID:             quiz.ID,
		TypeID:         quiz.TypeID,
		BadgeID:        quiz.BadgeID,
		ContinentID:    quiz.ContinentID,
		Country:        quiz.Country,
		Singular:       quiz.Singular,
		Name:           name,
		MaxScore:       quiz.MaxScore,
		Time:           quiz.Time,
		MapSVG:         quiz.MapSVG,
		ImageURL:       quiz.ImageURL,
		Plural:         plural,
		APIPath:        quiz.APIPath,
		Route:          quiz.Route,
		HasLeaderboard: quiz.HasLeaderboard,
		HasGrouping:    quiz.HasGrouping,
		HasFlags:       quiz.HasFlags,
		Enabled:        quiz.Enabled,
	}, nil
}

func GetQuiz(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	quiz, err := repo.GetQuiz(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(quiz)
}

func (s *Server) getQuizByRoute(writer http.ResponseWriter, request *http.Request) {
	quiz, err := repo.GetQuizByRoute(mux.Vars(request)["route"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	language := request.Header.Get("Content-Language")
	if language != "" && language != "en" {
		translatedQuiz, err := s.translateQuizDto(quiz, language)
		if err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(translatedQuiz)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(quiz)
}

func (s *Server) translateQuizDto(quiz repo.QuizDto, language string) (repo.QuizDto, error) {
	name, err := s.ts.TranslateText(language, quiz.Name)
	if err != nil {
		return repo.QuizDto{}, err
	}

	plural, err := s.ts.TranslateText(language, quiz.Plural)
	if err != nil {
		return repo.QuizDto{}, err
	}

	quizMap, err := s.getTranslatedMap(quiz, language)
	if err != nil {
		return repo.QuizDto{}, err
	}

	return repo.QuizDto{
		ID:             quiz.ID,
		TypeID:         quiz.TypeID,
		BadgeID:        quiz.BadgeID,
		ContinentID:    quiz.ContinentID,
		Country:        quiz.Country,
		Singular:       quiz.Singular,
		Name:           name,
		MaxScore:       quiz.MaxScore,
		Time:           quiz.Time,
		MapName:        quiz.MapName,
		Map:            quizMap,
		ImageURL:       quiz.ImageURL,
		Plural:         plural,
		APIPath:        quiz.APIPath,
		Route:          quiz.Route,
		HasLeaderboard: quiz.HasLeaderboard,
		HasGrouping:    quiz.HasGrouping,
		HasFlags:       quiz.HasFlags,
		Enabled:        quiz.Enabled,
	}, nil
}

func (s *Server) getTranslatedMap(quiz repo.QuizDto, language string) (repo.MapDto, error) {
	if quiz.MapName == "" {
		return quiz.Map, nil
	}

	translatedElements := make([]repo.MapElementDto, len(quiz.Map.Elements))
	for i, element := range quiz.Map.Elements {
		translatedElement, err := s.translateMapElementDto(element, language)
		if err != nil {
			return repo.MapDto{}, err
		}

		translatedElements[i] = translatedElement
	}

	return repo.MapDto{
		ID:        quiz.Map.ID,
		Key:       quiz.Map.Key,
		ClassName: quiz.Map.ClassName,
		Label:     quiz.Map.Label,
		ViewBox:   quiz.Map.ViewBox,
		Elements:  translatedElements,
	}, nil
}

func (s *Server) translateMapElementDto(element repo.MapElementDto, language string) (repo.MapElementDto, error) {
	name, err := s.ts.TranslateText(language, element.Name)
	if err != nil {
		return repo.MapElementDto{}, err
	}

	return repo.MapElementDto{
		EntryID:    element.EntryID,
		MapID:      element.MapID,
		Type:       element.Type,
		ID:         element.ID,
		Name:       name,
		D:          element.D,
		Points:     element.Points,
		X:          element.X,
		Y:          element.Y,
		Width:      element.Width,
		Height:     element.Height,
		Cx:         element.Cx,
		Cy:         element.Cy,
		R:          element.R,
		Transform:  element.Transform,
		XlinkHref:  element.XlinkHref,
		ClipPath:   element.ClipPath,
		ClipPathId: element.ClipPathId,
		X1:         element.X1,
		Y1:         element.Y1,
		X2:         element.X2,
		Y2:         element.Y2,
	}, nil
}

func CreateQuiz(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var newQuiz repo.CreateQuizDto
	err = json.Unmarshal(requestBody, &newQuiz)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	quiz, err := repo.CreateQuiz(newQuiz)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(quiz)
}

func UpdateQuiz(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var quiz repo.UpdateQuizDto
	err = json.Unmarshal(requestBody, &quiz)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if err = repo.UpdateQuiz(id, quiz); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteQuiz(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.DeleteQuiz(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
