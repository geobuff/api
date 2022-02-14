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
	ID       int                  `json:"id"`
	Type     string               `json:"type"`
	Question string               `json:"question"`
	Answers  []ManualTriviaAnswer `json:"answers"`
}

type CreateManualTriviaQuestionDto struct {
	TypeID      int                           `json:"typeId"`
	Question    string                        `json:"question"`
	Map         string                        `json:"map"`
	Highlighted string                        `json:"highlighted"`
	FlagCode    string                        `json:"flagCode"`
	ImageURL    string                        `json:"imageUrl"`
	Answers     []CreateManualTriviaAnswerDto `json:"answers"`
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

		answers, err := GetManualTriviaAnswers(question.ID)
		if err != nil {
			return nil, err
		}
		question.Answers = answers

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

func CreateManualTriviaQuestion(question CreateManualTriviaQuestionDto) error {
	statement := "INSERT INTO manualtriviaquestions (typeid, question, map, highlighted, flagcode, imageurl) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL).Scan(&id)
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

func DeleteManualTriviaQuestion(questionID int) error {
	Connection.QueryRow("DELETE FROM manualtriviaanswers WHERE manualTriviaQuestionId = $1;", questionID)
	var id int
	return Connection.QueryRow("DELETE FROM manualtriviaquestions WHERE id = $1 RETURNING id;", questionID).Scan(&id)
}

func GetManualTriviaQuestions(typeID int) ([]ManualTriviaQuestion, error) {
	rows, err := Connection.Query("SELECT * FROM manualtriviaquestions WHERE typeid = $1;", typeID)
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