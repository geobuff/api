package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ManualTriviaQuestion struct {
	ID          int          `json:"id"`
	TypeID      int          `json:"typeId"`
	CategoryID  int          `json:"categoryId"`
	Question    string       `json:"question"`
	Map         string       `json:"map"`
	Highlighted string       `json:"highlighted"`
	FlagCode    string       `json:"flagCode"`
	ImageURL    string       `json:"imageUrl"`
	Explainer   string       `json:"explainer"`
	LastUsed    sql.NullTime `json:"lastUsed"`
	QuizDate    sql.NullTime `json:"quizDate"`
	LastUpdated time.Time    `json:"lastUpdated"`
}

type ManualTriviaQuestionDto struct {
	ID          int                  `json:"id"`
	TypeID      int                  `json:"typeId"`
	Type        string               `json:"type"`
	CategoryId  int                  `json:"categoryId"`
	Category    string               `json:"category"`
	Question    string               `json:"question"`
	Map         string               `json:"map"`
	Highlighted string               `json:"highlighted"`
	FlagCode    string               `json:"flagCode"`
	ImageURL    string               `json:"imageUrl"`
	Explainer   string               `json:"explainer"`
	LastUsed    sql.NullTime         `json:"lastUsed"`
	QuizDate    sql.NullTime         `json:"quizDate"`
	LastUpdated time.Time            `json:"lastUpdated"`
	Answers     []ManualTriviaAnswer `json:"answers"`
}

type GetManualTriviaQuestionEntriesFilterParams struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TypeID     int    `json:"typeId"`
	CategoryID int    `json:"categoryId"`
	Question   string `json:"question"`
}

type CreateManualTriviaQuestionDto struct {
	TypeID      int                           `json:"typeId"`
	CategoryID  int                           `json:"categoryId"`
	Question    string                        `json:"question"`
	Map         string                        `json:"map"`
	Highlighted string                        `json:"highlighted"`
	FlagCode    string                        `json:"flagCode"`
	ImageURL    string                        `json:"imageUrl"`
	Explainer   string                        `json:"explainer"`
	QuizDate    sql.NullTime                  `json:"quizDate"`
	Answers     []CreateManualTriviaAnswerDto `json:"answers"`
}

type UpdateManualTriviaQuestionDto struct {
	TypeID      int                           `json:"typeId"`
	CategoryID  int                           `json:"categoryId"`
	Question    string                        `json:"question"`
	Map         string                        `json:"map"`
	Highlighted string                        `json:"highlighted"`
	FlagCode    string                        `json:"flagCode"`
	ImageURL    string                        `json:"imageUrl"`
	Explainer   string                        `json:"explainer"`
	QuizDate    sql.NullTime                  `json:"quizDate"`
	Answers     []UpdateManualTriviaAnswerDto `json:"answers"`
}

func GetAllManualTriviaQuestions(filterParams GetManualTriviaQuestionEntriesFilterParams) ([]ManualTriviaQuestionDto, error) {
	statement := "SELECT q.id, q.typeid, t.name, c.id, c.name, q.question, q.map, q.highlighted, q.flagcode, q.imageurl, q.lastused, q.quizDate, q.explainer, q.lastupdated FROM manualtriviaquestions q JOIN triviaquestiontype t ON t.id = q.typeid JOIN triviaquestioncategory c ON c.id = q.categoryid WHERE q.question ILIKE '%' || $1 || '%' " + getTypeFilter(filterParams.TypeID) + getCategoryFilter(filterParams.CategoryID) + " ORDER BY q.lastupdated DESC LIMIT $2 OFFSET $3;"
	rows, err := Connection.Query(statement, filterParams.Question, filterParams.Limit, filterParams.Page*filterParams.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []ManualTriviaQuestionDto{}
	for rows.Next() {
		var question ManualTriviaQuestionDto
		if err = rows.Scan(&question.ID, &question.TypeID, &question.Type, &question.CategoryId, &question.Category, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL, &question.LastUsed, &question.QuizDate, &question.Explainer, &question.LastUpdated); err != nil {
			return nil, err
		}

		answers, err := GetManualTriviaAnswers(question.ID)
		if err != nil {
			return nil, err
		}
		question.Answers = answers

		questions = append(questions, question)
	}
	return questions, rows.Err()
}

func getTypeFilter(typeID int) string {
	if typeID == 0 {
		return ""
	}
	return fmt.Sprintf(" AND t.id = %d ", typeID)
}

func getCategoryFilter(categoryID int) string {
	if categoryID == 0 {
		return ""
	}
	return fmt.Sprintf(" AND c.id = %d ", categoryID)
}

var GetFirstManualTriviaQuestionID = func(filterParams GetManualTriviaQuestionEntriesFilterParams) (int, error) {
	statement := "SELECT q.id FROM manualtriviaquestions q JOIN triviaquestiontype t ON t.id = q.typeid JOIN triviaquestioncategory c ON c.id = q.categoryid WHERE q.question ILIKE '%' || $1 || '%' " + getTypeFilter(filterParams.TypeID) + getCategoryFilter(filterParams.CategoryID) + " ORDER BY q.lastupdated DESC LIMIT 1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, filterParams.Question, (filterParams.Page+1)*filterParams.Limit).Scan(&id)
	return id, err
}

func CreateManualTriviaQuestion(question CreateManualTriviaQuestionDto) error {
	statement := "INSERT INTO manualtriviaquestions (typeid, categoryid, question, map, highlighted, flagcode, imageurl, quizDate, explainer, lastupdated) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.TypeID, question.CategoryID, strings.TrimSpace(question.Question), question.Map, question.Highlighted, question.FlagCode, question.ImageURL, question.QuizDate, question.Explainer, time.Now()).Scan(&id)
	if err != nil {
		return err
	}

	for _, val := range question.Answers {
		if err = CreateManualTriviaAnswer(id, val); err != nil {
			return err
		}
	}
	return nil
}

func ValidateCreateQuestion(question CreateManualTriviaQuestionDto) error {
	var id int
	err := Connection.QueryRow("SELECT id FROM manualtriviaquestions WHERE question ILIKE '%' || $1 || '%';", question.Question).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil {
		return errors.New(fmt.Sprintf("Error! Question with text %s already exists.", question.Question))
	}

	visited := make(map[string]bool)
	for _, answer := range question.Answers {
		lower := strings.ToLower(answer.Text)
		if visited[lower] {
			return errors.New(fmt.Sprintf("Error! Duplicate answers with text %s.", answer.Text))
		}
		visited[lower] = true
	}
	return nil
}

func UpdateManualTriviaQuestion(questionID int, question UpdateManualTriviaQuestionDto) error {
	statement := "UPDATE manualtriviaquestions SET typeid = $2, categoryid = $3, question = $4, map = $5, highlighted = $6, flagcode = $7, imageurl = $8, quizDate = $9, explainer = $10, lastupdated = $11 WHERE id = $1 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, questionID, question.TypeID, question.CategoryID, strings.TrimSpace(question.Question), question.Map, question.Highlighted, question.FlagCode, question.ImageURL, question.QuizDate, question.Explainer, time.Now()).Scan(&id)
	if err != nil {
		return err
	}

	for _, val := range question.Answers {
		if err = UpdateManualTriviaAnswer(val); err != nil {
			return err
		}
	}
	return nil
}

func ValidateUpdateQuestion(question UpdateManualTriviaQuestionDto) error {
	visited := make(map[string]bool)
	for _, answer := range question.Answers {
		lower := strings.ToLower(answer.Text)
		if visited[lower] {
			return errors.New(fmt.Sprintf("Error! Duplicate answers with text %s.", answer.Text))
		}
		visited[lower] = true
	}
	return nil
}

func DeleteManualTriviaQuestion(questionID int) error {
	Connection.QueryRow("DELETE FROM manualtriviaanswers WHERE manualTriviaQuestionId = $1;", questionID)
	var id int
	return Connection.QueryRow("DELETE FROM manualtriviaquestions WHERE id = $1 RETURNING id;", questionID).Scan(&id)
}

func GetManualTriviaQuestions(typeID int, lastUsedMax string) ([]ManualTriviaQuestion, error) {
	statement := "SELECT * FROM manualtriviaquestions WHERE typeid = $1 AND quizdate IS null AND (lastUsed IS null OR lastUsed < $2);"
	rows, err := Connection.Query(statement, typeID, lastUsedMax)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []ManualTriviaQuestion{}
	for rows.Next() {
		var question ManualTriviaQuestion
		if err = rows.Scan(&question.ID, &question.TypeID, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL, &question.LastUsed, &question.QuizDate, &question.Explainer, &question.LastUpdated); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

func UpdateManualTriviaQuestionLastUsed(questionID int) error {
	today := time.Now().Format("2006-01-02")
	statement := "UPDATE manualtriviaquestions SET lastUsed = $2 WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionID, today).Scan(&id)
}

func GetTodaysManualTriviaQuestions() ([]ManualTriviaQuestion, error) {
	today := time.Now().Format("2006-01-02")
	rows, err := Connection.Query("SELECT * FROM manualtriviaquestions WHERE quizDate = $1;", today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []ManualTriviaQuestion{}
	for rows.Next() {
		var question ManualTriviaQuestion
		if err = rows.Scan(&question.ID, &question.TypeID, &question.CategoryID, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL, &question.LastUsed, &question.QuizDate, &question.Explainer, &question.LastUpdated); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}
