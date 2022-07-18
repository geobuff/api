package repo

import (
	"database/sql"
	"math/rand"
)

type CommunityQuizQuestion struct {
	ID                 int    `json:"id"`
	CommunityQuizID    int    `json:"communityQuizId"`
	TypeID             int    `json:"typeId"`
	Question           string `json:"question"`
	Map                string `json:"map"`
	Highlighted        string `json:"highlighted"`
	FlagCode           string `json:"flagCode"`
	ImageUrl           string `json:"imageUrl"`
	ImageAttributeName string `json:"imageAttributeName"`
	ImageAttributeURL  string `json:"imageAttributeUrl"`
	ImageWidth         int    `json:"imageWidth"`
	ImageHeight        int    `json:"imageHeight"`
	ImageAlt           string `json:"imageAlt"`
	Explainer          string `json:"explainer"`
}

type GetCommunityQuizQuestionDto struct {
	ID                 int                         `json:"id"`
	TypeID             int                         `json:"typeId"`
	Type               string                      `json:"type"`
	Question           string                      `json:"question"`
	Map                string                      `json:"map"`
	Highlighted        string                      `json:"highlighted"`
	FlagCode           string                      `json:"flagCode"`
	ImageUrl           string                      `json:"imageUrl"`
	ImageAttributeName string                      `json:"imageAttributeName"`
	ImageAttributeURL  string                      `json:"imageAttributeUrl"`
	ImageWidth         int                         `json:"imageWidth"`
	ImageHeight        int                         `json:"imageHeight"`
	ImageAlt           string                      `json:"imageAlt"`
	Explainer          string                      `json:"explainer"`
	Answers            []GetCommunityQuizAnswerDto `json:"answers"`
}

type CreateCommunityQuizQuestionDto struct {
	TypeID             int                            `json:"typeId"`
	Question           string                         `json:"question"`
	Map                string                         `json:"map"`
	Highlighted        string                         `json:"highlighted"`
	FlagCode           string                         `json:"flagCode"`
	ImageUrl           string                         `json:"imageUrl"`
	ImageAttributeName string                         `json:"imageAttributeName"`
	ImageAttributeURL  string                         `json:"imageAttributeUrl"`
	ImageWidth         int                            `json:"imageWidth"`
	ImageHeight        int                            `json:"imageHeight"`
	ImageAlt           string                         `json:"imageAlt"`
	Explainer          string                         `json:"explainer"`
	Answers            []CreateCommunityQuizAnswerDto `json:"answers"`
}

type UpdateCommunityQuizQuestionDto struct {
	ID                 sql.NullInt64                  `json:"id"`
	TypeID             int                            `json:"typeId"`
	Question           string                         `json:"question"`
	Map                string                         `json:"map"`
	Highlighted        string                         `json:"highlighted"`
	FlagCode           string                         `json:"flagCode"`
	ImageUrl           string                         `json:"imageUrl"`
	ImageAttributeName string                         `json:"imageAttributeName"`
	ImageAttributeURL  string                         `json:"imageAttributeUrl"`
	ImageWidth         int                            `json:"imageWidth"`
	ImageHeight        int                            `json:"imageHeight"`
	ImageAlt           string                         `json:"imageAlt"`
	Explainer          string                         `json:"explainer"`
	Answers            []CreateCommunityQuizAnswerDto `json:"answers"`
}

func InsertCommunityQuizQuestion(quizID int, question CreateCommunityQuizQuestionDto) (int, error) {
	statement := "INSERT INTO communityquizquestions (communityquizid, typeid, question, map, highlighted, flagcode, imageurl, imageAttributeName, imageAttributeUrl, imageWidth, imageHeight, imageAlt, explainer) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, quizID, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageUrl, question.ImageAttributeName, question.ImageAttributeURL, question.ImageWidth, question.ImageHeight, question.ImageAlt, question.Explainer).Scan(&id)
	return id, err
}

func UpdateCommunityQuizQuestion(questionID int, question UpdateCommunityQuizQuestionDto) error {
	statement := "UPDATE communityquizquestions SET typeid = $1, question = $2, map = $3, highlighted = $4, flagcode = $5, imageurl = $6, imageAttributeName = $7, imageAttributeUrl = $8, imageWidth = $9, imageHeight = $10, imageAlt = $11, explainer = $12 WHERE id = $13 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageUrl, question.ImageAttributeName, question.ImageAttributeURL, question.ImageWidth, question.ImageHeight, question.ImageAlt, question.Explainer, questionID).Scan(&id)
}

func GetCommunityQuizQuestionIds(quizID int) ([]int, error) {
	statement := "SELECT id FROM communityquizquestions WHERE communityquizid = $1;"
	rows, err := Connection.Query(statement, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids = []int{}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func GetCommunityQuizQuestions(quizID int) ([]GetCommunityQuizQuestionDto, error) {
	statement := "SELECT q.id, q.typeid, t.name, q.question, q.map, q.highlighted, q.flagcode, q.imageurl, q.imageAttributeName, q.imageAttributeUrl, q.imageWidth, q.imageHeight, q.imageAlt, q.explainer FROM communityquizquestions q JOIN triviaQuestionType t ON t.id = q.typeid WHERE communityquizid = $1;"
	rows, err := Connection.Query(statement, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []GetCommunityQuizQuestionDto{}
	for rows.Next() {
		var question GetCommunityQuizQuestionDto
		if err = rows.Scan(&question.ID, &question.TypeID, &question.Type, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageUrl, &question.ImageAttributeName, &question.ImageAttributeURL, &question.ImageWidth, &question.ImageHeight, &question.ImageAlt, &question.Explainer); err != nil {
			return nil, err
		}

		answers, err := GetCommunityQuizAnswers(question.ID)
		if err != nil {
			return nil, err
		}

		if len(answers) > 2 {
			rand.Shuffle(len(answers), func(i, j int) {
				answers[i], answers[j] = answers[j], answers[i]
			})
		}
		question.Answers = answers

		questions = append(questions, question)
	}
	return questions, rows.Err()
}

func DeleteCommunityQuizQuestion(questionID int) error {
	statement := "DELETE FROM communityquizquestions WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionID).Scan(&id)
}
