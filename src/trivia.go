package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/utils"
	"github.com/gorilla/mux"
)

type GetTriviaDto struct {
	Trivia  []repo.Trivia `json:"trivia"`
	HasMore bool          `json:"hasMore"`
}

func GenerateTrivia(writer http.ResponseWriter, request *http.Request) {
	err := repo.CreateTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func RegenerateTrivia(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err := repo.RegenerateTrivia(mux.Vars(request)["date"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteTrivia(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err := repo.DeleteTriviaByDate(mux.Vars(request)["date"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteOldTrivia(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	newTriviaCount, err := strconv.Atoi(mux.Vars(request)["newTriviaCount"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.DeleteOldTrivia(newTriviaCount)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func GetAllTrivia(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var filter repo.GetTriviaFilter
	err = json.Unmarshal(requestBody, &filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	trivia, err := repo.GetAllTrivia(filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	language := request.Header.Get("Content-Language")
	if language != "" && language != "en" {
		translatedTrivia := make([]repo.Trivia, len(trivia))
		for index, quiz := range trivia {
			name, err := utils.TranslateText(language, quiz.Name)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			translatedTrivia[index] = repo.Trivia{
				ID:       quiz.ID,
				Name:     name,
				Date:     quiz.Date,
				MaxScore: quiz.MaxScore,
			}
		}

		switch _, err := repo.GetFirstTriviaID(filter); err {
		case sql.ErrNoRows:
			entriesDto := GetTriviaDto{translatedTrivia, false}
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(entriesDto)
		case nil:
			entriesDto := GetTriviaDto{translatedTrivia, true}
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(entriesDto)
		default:
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		}
		return
	}

	switch _, err := repo.GetFirstTriviaID(filter); err {
	case sql.ErrNoRows:
		entriesDto := GetTriviaDto{trivia, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := GetTriviaDto{trivia, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

func GetTriviaByDate(writer http.ResponseWriter, request *http.Request) {
	trivia, err := repo.GetTrivia(mux.Vars(request)["date"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	language := request.Header.Get("Content-Language")
	if language != "" && language != "en" {
		name, err := utils.TranslateText(language, trivia.Name)
		if err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
			return
		}

		translatedTrivia := repo.TriviaDto{
			ID:        trivia.ID,
			Name:      name,
			MaxScore:  trivia.MaxScore,
			Questions: make([]repo.QuestionDto, len(trivia.Questions)),
		}

		for index, question := range trivia.Questions {
			questionValue, err := utils.TranslateText(language, question.Question)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			imageAlt, err := utils.TranslateText(language, question.ImageAlt)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			explainer, err := utils.TranslateText(language, question.Explainer)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			answers := make([]repo.AnswerDto, len(question.Answers))
			for index, answer := range question.Answers {
				text, err := utils.TranslateText(language, answer.Text)
				if err != nil {
					http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
					return
				}

				answers[index] = repo.AnswerDto{
					Text:      text,
					IsCorrect: answer.IsCorrect,
					FlagCode:  answer.FlagCode,
					FlagUrl:   answer.FlagUrl,
				}
			}

			translatedTrivia.Questions[index] = repo.QuestionDto{
				ID:                 question.ID,
				Type:               question.Type,
				Question:           questionValue,
				MapName:            question.MapName,
				Map:                question.Map,
				Highlighted:        question.Highlighted,
				FlagCode:           question.FlagCode,
				FlagUrl:            question.FlagUrl,
				ImageURL:           question.ImageURL,
				ImageAttributeName: question.ImageAttributeName,
				ImageAttributeURL:  question.ImageAttributeURL,
				ImageWidth:         question.ImageWidth,
				ImageHeight:        question.ImageHeight,
				ImageAlt:           imageAlt,
				Explainer:          explainer,
				Answers:            answers,
			}
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(translatedTrivia)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(trivia)
}
