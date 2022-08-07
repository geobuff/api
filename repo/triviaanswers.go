package repo

import (
	"database/sql"
	"math/rand"
)

type TriviaAnswer struct {
	ID               int    `json:"id"`
	TriviaQuestionID int    `json:"triviaQuestionId"`
	Text             string `json:"text"`
	IsCorrect        bool   `json:"isCorrect"`
	FlagCode         string `json:"flagCode"`
}

type AnswerDto struct {
	Text      string         `json:"text"`
	IsCorrect bool           `json:"isCorrect"`
	FlagCode  string         `json:"flagCode"`
	FlagUrl   sql.NullString `json:"flagUrl"`
}

func GetTriviaAnswers(triviaQuestionId int) ([]AnswerDto, error) {
	rows, err := Connection.Query("SELECT a.text, a.isCorrect, a.flagCode, f.url FROM triviaAnswers a LEFT JOIN flagentries f ON f.code = a.flagcode WHERE triviaQuestionId = $1;", triviaQuestionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers = []AnswerDto{}
	for rows.Next() {
		var answer AnswerDto
		if err = rows.Scan(&answer.Text, &answer.IsCorrect, &answer.FlagCode, &answer.FlagUrl); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	if len(answers) > 2 {
		rand.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
		})
	}

	return answers, nil
}

func CreateTriviaAnswer(answer TriviaAnswer) error {
	statement := "INSERT INTO triviaAnswers (triviaQuestionId, text, isCorrect, flagCode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, answer.TriviaQuestionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}

func DeleteTriviaAnswers(triviaQuestionId int) error {
	statement := "DELETE FROM triviaAnswers WHERE triviaQuestionId = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, triviaQuestionId).Scan(&id)
}
