package repo

import (
	"database/sql"
	"time"
)

type Trivia struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	MaxScore int       `json:"maxScore"`
}

type TriviaDto struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	MaxScore  int           `json:"maxScore"`
	Questions []QuestionDto `json:"questions"`
}

type GetTriviaFilter struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

func GetAllTrivia(filter GetTriviaFilter) ([]Trivia, error) {
	statement := "SELECT * FROM trivia WHERE name ILIKE '%' || $1 || '%' ORDER BY date DESC LIMIT $2 OFFSET $3;"
	rows, err := Connection.Query(statement, filter.Filter, filter.Limit, filter.Page*filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trivia = []Trivia{}
	for rows.Next() {
		var quiz Trivia
		if err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.Date, &quiz.MaxScore); err != nil {
			return nil, err
		}
		trivia = append(trivia, quiz)
	}
	return trivia, rows.Err()
}

func GetFirstTriviaID(filter GetTriviaFilter) (int, error) {
	statement := "SELECT id FROM trivia WHERE name ILIKE '%' || $1 || '%' ORDER BY date DESC LIMIT 1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, filter.Filter, (filter.Page+1)*filter.Page).Scan(&id)
	return id, err
}

func GetTrivia(date string) (*TriviaDto, error) {
	var result TriviaDto
	err := Connection.QueryRow("SELECT id, name, maxscore from trivia WHERE date = $1;", date).Scan(&result.ID, &result.Name, &result.MaxScore)
	if err != nil {
		return nil, err
	}

	questions, err := GetTriviaQuestions(result.ID)
	if err != nil {
		return nil, err
	}

	result.Questions = questions
	return &result, nil
}

func DeleteTriviaByDate(dateString string) error {
	trivia, err := GetTrivia(dateString)
	if err != nil {
		return err
	}

	return deleteTrivia(trivia)
}

func deleteTrivia(trivia *TriviaDto) error {
	if err := ClearTriviaPlayTriviaId(trivia.ID); err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, question := range trivia.Questions {
		if err := DeleteTriviaAnswers(question.ID); err != nil && err != sql.ErrNoRows {
			return err
		}

		if err := DeleteTriviaQuestion(question.ID); err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	var id int
	return Connection.QueryRow("DELETE FROM trivia WHERE id = $1 RETURNING id;", trivia.ID).Scan(&id)
}

func DeleteOldTrivia(newTriviaCount int) error {
	dates, err := getOldTriviaDates(newTriviaCount)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	for _, date := range dates {
		formattedDate := date.Format("2006-01-02")
		trivia, err := GetTrivia(formattedDate)
		if err != nil {
			return err
		}

		if err = deleteTrivia(trivia); err != nil {
			return err
		}
	}
	return err
}

func getOldTriviaDates(newTriviaCount int) ([]time.Time, error) {
	rows, err := Connection.Query("SELECT date FROM trivia WHERE date < $1;", time.Now().AddDate(0, 0, 0-newTriviaCount))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates = []time.Time{}
	for rows.Next() {
		var date time.Time
		if err = rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return dates, rows.Err()
}
