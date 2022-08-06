package repo

import (
	"database/sql"
	"math/rand"
)

type TriviaQuestion struct {
	ID                 int    `json:"id"`
	TriviaId           int    `json:"triviaId"`
	TypeID             int    `json:"typeId"`
	Question           string `json:"question"`
	Map                string `json:"map"`
	Highlighted        string `json:"highlighted"`
	FlagCode           string `json:"flagCode"`
	ImageURL           string `json:"imageUrl"`
	ImageAttributeName string `json:"imageAttributeName"`
	ImageAttributeURL  string `json:"ImageAttributeUrl"`
	ImageWidth         int    `json:"imageWidth"`
	ImageHeight        int    `json:"imageHeight"`
	ImageAlt           string `json:"imageAlt"`
	Explainer          string `json:"explainer"`
}

type QuestionDto struct {
	ID                 int            `json:"id"`
	Type               string         `json:"type"`
	Question           string         `json:"question"`
	MapName            string         `json:"mapName"`
	Map                MapDto         `json:"map"`
	Highlighted        string         `json:"highlighted"`
	FlagCode           string         `json:"flagCode"`
	FlagUrl            sql.NullString `json:"flagUrl"`
	ImageURL           string         `json:"imageUrl"`
	ImageAttributeName string         `json:"imageAttributeName"`
	ImageAttributeURL  string         `json:"imageAttributeUrl"`
	ImageWidth         int            `json:"imageWidth"`
	ImageHeight        int            `json:"imageHeight"`
	ImageAlt           string         `json:"imageAlt"`
	Explainer          string         `json:"explainer"`
	Answers            []AnswerDto    `json:"answers"`
}

func GetTriviaQuestions(triviaId int) ([]QuestionDto, error) {
	rows, err := Connection.Query("SELECT q.id, t.name, q.question, q.map, q.highlighted, q.flagCode, f.url, q.imageUrl, q.imageAttributeName, q.imageAttributeUrl, q.imageWidth, q.imageHeight, q.imageAlt, q.explainer FROM triviaQuestions q JOIN triviaQuestionType t ON t.id = q.typeId LEFT JOIN flagEntries f ON f.code = q.flagCode WHERE q.triviaId = $1;", triviaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []QuestionDto{}
	for rows.Next() {
		var question QuestionDto
		if err = rows.Scan(&question.ID, &question.Type, &question.Question, &question.MapName, &question.Highlighted, &question.FlagCode, &question.FlagUrl, &question.ImageURL, &question.ImageAttributeName, &question.ImageAttributeURL, &question.ImageWidth, &question.ImageHeight, &question.ImageAlt, &question.Explainer); err != nil {
			return nil, err
		}

		if question.MapName != "" {
			svgMap, err := GetMap(question.MapName)
			if err != nil {
				return nil, err
			}
			question.Map = svgMap
		}

		answers, err := GetTriviaAnswers(question.ID)
		if err != nil {
			return nil, err
		}

		question.Answers = answers
		questions = append(questions, question)
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return questions, nil
}

func CreateTriviaQuestion(question TriviaQuestion) (int, error) {
	statement := "INSERT INTO triviaQuestions (triviaId, typeId, question, map, highlighted, flagCode, imageUrl, imageAttributeName, imageAttributeUrl, imageWidth, imageHeight, imageAlt, explainer) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.TriviaId, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL, question.ImageAttributeName, question.ImageAttributeURL, question.ImageWidth, question.ImageHeight, question.ImageAlt, question.Explainer).Scan(&id)
	return id, err
}

func DeleteTriviaQuestion(questionId int) error {
	statement := "DELETE FROM triviaQuestions WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionId).Scan(&id)
}
