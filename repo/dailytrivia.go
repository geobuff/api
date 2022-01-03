package repo

import (
	"database/sql"
	"errors"
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
	Highlighted   string `json:"hightlighted"`
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

type DailyTriviaDto struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Questions []QuestionDto `json:"questions"`
}

type QuestionDto struct {
	ID          int         `json:"id"`
	Type        string      `json:"type"`
	Question    string      `json:"question"`
	Map         string      `json:"map"`
	Highlighted string      `json:"highlighted"`
	FlagCode    string      `json:"flagCode"`
	ImageURL    string      `json:"imageUrl"`
	Answers     []AnswerDto `json:"answers"`
}

type AnswerDto struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

func CreateDailyTrivia() error {
	date := time.Now()
	var id int
	if err := Connection.QueryRow("SELECT id FROM dailyTrivia WHERE date = $1", date).Scan(&id); err != sql.ErrNoRows {
		return errors.New("daily trivia for current date already created")
	}

	_, month, day := date.Date()
	statement := "INSERT INTO dailyTrivia (name, date) VALUES ($1, $2) RETURNING id;"
	err := Connection.QueryRow(statement, fmt.Sprintf("Daily Trivia - %s %d", month, day), date).Scan(&id)
	if err != nil {
		return err
	}
	return generateQuestions(id)
}

func generateQuestions(dailyTriviaId int) error {
	err := whatCountry(dailyTriviaId)
	if err != nil {
		return err
	}

	err = whatCapital(dailyTriviaId)
	if err != nil {
		return err
	}

	err = whatFlag(dailyTriviaId)
	if err != nil {
		return err
	}

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

func whatCountry(dailyTriviaId int) error {
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
		Question:      "Which country is highlighted above?",
		Map:           "WorldCountries",
		Highlighted:   country.SVGName,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  country.SVGName,
		IsCorrect:             true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max + 1)
		country = countries[index]
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  country.SVGName,
			IsCorrect:             false,
		}

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatCapital(dailyTriviaId int) error {
	capitals := mapping.Mappings["world-capitals"]
	rand.Seed(time.Now().UnixNano())
	max := len(capitals)
	index := rand.Intn(max + 1)
	capital := capitals[index]
	capitals = append(capitals[:index], capitals[index+1:]...)
	max = max - 1
	countryName := getCountryName(capital.Code)

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "map",
		Question:      fmt.Sprintf("The capital city of %s is _______", countryName),
		Map:           "WorldCapitals",
		Highlighted:   capital.SVGName,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  capital.SVGName,
		IsCorrect:             true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max + 1)
		capital := capitals[index]
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  capital.SVGName,
			IsCorrect:             false,
		}

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		capitals = append(capitals[:index], capitals[index+1:]...)
		max = max - 1
	}

	return nil
}

func getCountryName(code string) string {
	for _, value := range mapping.Mappings["world-countries"] {
		if value.Code == code {
			return value.SVGName
		}
	}
	return ""
}

func whatFlag(dailyTriviaId int) error {
	countries := mapping.Mappings["world-countries"]
	rand.Seed(time.Now().UnixNano())
	max := len(countries)
	index := rand.Intn(max + 1)
	country := countries[index]
	countries = append(countries[:index], countries[index+1:]...)
	max = max - 1

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "flag",
		Question:      "Which country has this flag?",
		FlagCode:      country.Code,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  country.SVGName,
		IsCorrect:             true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max + 1)
		country = countries[index]
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  country.SVGName,
			IsCorrect:             false,
		}

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}
func GetAllDailyTrivia() ([]DailyTrivia, error) {
	rows, err := Connection.Query("SELECT * FROM dailyTrivia LIMIT 20;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trivia = []DailyTrivia{}
	for rows.Next() {
		var quiz DailyTrivia
		if err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.Date); err != nil {
			return nil, err
		}
		trivia = append(trivia, quiz)
	}
	return trivia, rows.Err()
}

func GetDailyTrivia(date string) (*DailyTriviaDto, error) {
	var result DailyTriviaDto
	err := Connection.QueryRow("SELECT id, name from dailyTrivia WHERE date = $1;", date).Scan(&result.ID, &result.Name)
	if err != nil {
		return nil, err
	}

	rows, err := Connection.Query("SELECT id, type, question, map, highlighted, flagCode, imageUrl FROM dailyTriviaQuestions WHERE dailyTriviaId = $1;", result.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []QuestionDto{}
	for rows.Next() {
		var question QuestionDto
		if err = rows.Scan(&question.ID, &question.Type, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL); err != nil {
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
	result.Questions = questions

	return &result, rows.Err()
}