package repo

type ManualTriviaQuestion struct {
	ID          int    `json:"id"`
	TypeID      int    `json:"typeId"`
	Question    string `json:"question"`
	Map         string `json:"map"`
	Highlighted string `json:"highlighted"`
	FlagCode    string `json:"flagCode"`
	ImageURL    string `json:"imageUrl"`
}

func GetAllManualTriviaQuestions(limit, page int) ([]ManualTriviaQuestion, error) {
	rows, err := Connection.Query("SELECT * FROM manualtriviaquestions LIMIT $1 OFFSET $2;", limit, limit*page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []ManualTriviaQuestion{}
	for rows.Next() {
		var question ManualTriviaQuestion
		if err = rows.Scan(&question.ID, &question.TypeID, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

var GetFirstManualTriviaQuestionID = func(offset int) (int, error) {
	statement := "SELECT id FROM manualtriviaquestions LIMIT 1 OFFSET $1;"
	var id int
	err := Connection.QueryRow(statement, offset).Scan(&id)
	return id, err
}
