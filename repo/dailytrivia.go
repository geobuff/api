package repo

import (
	"fmt"
	"time"
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

var questions = []DailyTriviaQuestion{
	{
		DailyTriviaId: 1,
		Type:          "text",
		Question:      "If I’m visiting the ancient city of Petra, which country am I in?",
	},
	{
		DailyTriviaId: 1,
		Type:          "flag",
		Question:      "What country does this flag belong to?",
		FlagCode:      "nz",
	},
	{
		DailyTriviaId: 1,
		Type:          "map",
		Question:      "What Country is this?",
		Map:           "NewZealandRegions",
	},
	{
		DailyTriviaId: 1,
		Type:          "text",
		Question:      "Is Australia Real?",
	},
}

var answers = []DailyTriviaAnswer{
	{
		DailyTriviaQuestionID: 1,
		Text:                  "Peru",
		FlagCode:              "pe",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 1,
		Text:                  "United Arab Emirates",
		FlagCode:              "ae",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 1,
		Text:                  "Jordan",
		FlagCode:              "jo",
		IsCorrect:             true,
	},
	{
		DailyTriviaQuestionID: 1,
		Text:                  "Jeff Bezos",
		FlagCode:              "us",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 2,
		Text:                  "Peru",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 2,
		Text:                  "United Arab Emirates",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 2,
		Text:                  "Jordan",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 2,
		Text:                  "New Zealand",
		IsCorrect:             true,
	},
	{
		DailyTriviaQuestionID: 3,
		Text:                  "Peru",
		FlagCode:              "pe",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 3,
		Text:                  "United Arab Emirates",
		FlagCode:              "ae",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 3,
		Text:                  "New Zealand",
		FlagCode:              "nz",
		IsCorrect:             true,
	},
	{
		DailyTriviaQuestionID: 3,
		Text:                  "Jeff Bezos",
		FlagCode:              "us",
		IsCorrect:             false,
	},
	{
		DailyTriviaQuestionID: 4,
		Text:                  "True",
		IsCorrect:             true,
	},
	{
		DailyTriviaQuestionID: 4,
		Text:                  "False",
		IsCorrect:             false,
	},
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

	for _, value := range questions {
		statement = "INSERT INTO dailyTriviaQuestions (dailyTriviaId, type, question, map, flagCode, imageUrl) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
		var id int
		err := Connection.QueryRow(statement, value.DailyTriviaId, value.Type, value.Question, value.Map, value.FlagCode, value.ImageURL).Scan(&id)
		if err != nil {
			return err
		}
	}

	for _, value := range answers {
		statement = "INSERT INTO dailyTriviaAnswers (dailyTriviaQuestionId, text, isCorrect, flagCode) VALUES ($1, $2, $3, $4) RETURNING id;"
		var id int
		err := Connection.QueryRow(statement, value.DailyTriviaQuestionID, value.Text, value.IsCorrect, value.FlagCode).Scan(&id)
		if err != nil {
			return err
		}
	}

	return err
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
