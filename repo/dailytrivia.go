package repo

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/geobuff/mapping"
)

type DailyTrivia struct {
	ID   int       `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type DailyTriviaQuestion struct {
	ID            int    `json:"id"`
	DailyTriviaId int    `json:"dailyTriviaId"`
	Type          string `json:"type"`
	Question      string `json:"question"`
	Map           string `json:"map"`
	Highlighted   string `json:"hightlighted`
	FlagCode      string `json:"flagCode"`
	ImageURL      string `json:"imageUrl"`
}

type DailyTriviaAnswer struct {
	ID                    int    `json:"id"`
	DailyTriviaQuestionID int    `json:"dailyTriviaQuestionId"`
	Text                  string `json:"text"`
	IsCorrect             bool   `json:"isCorrect"`
	FlagCode              string `json:"flagCode"`
}

type QuestionDto struct {
	ID       int         `json:"id"`
	Type     string      `json:"type"`
	Question string      `json:"question"`
	Map      string      `json:"map"`
	FlagCode string      `json:"flagCode"`
	ImageURL string      `json:"imageUrl"`
	Answers  []AnswerDto `json:"answers"`
}

type AnswerDto struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

func CreateDailyTrivia() error {
	date := time.Now()
	year, month, day := date.Date()
	statement := "INSERT INTO dailyTrivia (name, date) VALUES ($1, $2) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, fmt.Sprintf("Daily Trivia %d-%d-%d", day, month, year), date).Scan(&id)
	if err != nil {
		return err
	}

	return generateQuestions(id)
}

func generateQuestions(dailyTriviaId int) error {
	questionId, err := whatCountry(dailyTriviaId)
	if err != nil {
		return err
	}
	questionId = questionId + 1
	return nil
}

func createQuestion(question DailyTriviaQuestion) (int, error) {
	statement := "INSERT INTO dailyTriviaQuestions (dailyTriviaId, type, question, map, highlighted, flagCode, imageUrl) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.DailyTriviaId, question.Type, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL).Scan(&id)
	return id, err
}

func createAnswer(answer DailyTriviaAnswer) error {
	statement := "INSERT INTO dailyTriviaAnswers (dailyTriviaQuestionId, text, isCorrect, flagCode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, answer.DailyTriviaQuestionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}

func whatCountry(dailyTriviaId int) (int, error) {
	countries := mapping.Mappings["world-countries"]
	rand.Seed(time.Now().UnixNano())
	max := len(countries)
	index := rand.Intn(max + 1)
	country := countries[index]
	countries = append(countries[:index], countries[index+1:]...)
	max = max - 1

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "map",
		Question:      "What country is highlighted below?",
		Map:           "WorldCountries",
		Highlighted:   country.SVGName,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return 0, err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  country.SVGName,
		IsCorrect:             true,
		FlagCode:              country.Code,
	}
	err = createAnswer(answer)
	if err != nil {
		return 0, err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max + 1)
		country = countries[index]
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  country.SVGName,
			IsCorrect:             false,
			FlagCode:              country.Code,
		}

		err = createAnswer(answer)
		if err != nil {
			return 0, err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return questionId, nil
}

func GetDailyTrivia(date string) ([]QuestionDto, error) {
	var id int
	err := Connection.QueryRow("SELECT id from dailyTrivia WHERE date = $1;", date).Scan(&id)
	if err != nil {
		return nil, err
	}

	rows, err := Connection.Query("SELECT id, type, question, map, flagCode, imageUrl FROM dailyTriviaQuestions WHERE dailyTriviaId = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []QuestionDto{}
	for rows.Next() {
		var question QuestionDto
		if err = rows.Scan(&question.ID, &question.Type, &question.Question, &question.Map, &question.FlagCode, &question.ImageURL); err != nil {
			return nil, err
		}

		rows, err := Connection.Query("SELECT text, isCorrect, flagCode FROM dailyTriviaAnswers WHERE dailyTriviaQuestionId = $1;", question.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var answers = []AnswerDto{}
		for rows.Next() {
			var answer AnswerDto
			if err = rows.Scan(&answer.Text, &answer.IsCorrect, &answer.FlagCode); err != nil {
				return nil, err
			}
			answers = append(answers, answer)
		}
		question.Answers = answers

		questions = append(questions, question)
	}

	return questions, rows.Err()
}
