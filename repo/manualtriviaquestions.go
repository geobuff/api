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

type ManualTriviaQuestionDto struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Question string `json:"question"`
}

func GetAllManualTriviaQuestions(limit, offset int) ([]ManualTriviaQuestionDto, error) {
	rows, err := Connection.Query("SELECT q.id, t.name, q.question FROM manualtriviaquestions q JOIN triviaquestiontype t ON t.id = q.typeid LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []ManualTriviaQuestionDto{}
	for rows.Next() {
		var question ManualTriviaQuestionDto
		if err = rows.Scan(&question.ID, &question.Type, &question.Question); err != nil {
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

func DeleteManualTriviaQuestion(questionID int) error {
	Connection.QueryRow("DELETE FROM manualtriviaanswers WHERE manualTriviaQuestionId = $1;", questionID)
	var id int
	return Connection.QueryRow("DELETE FROM manualtriviaquestions WHERE id = $1 RETURNING id;", questionID).Scan(&id)
}
