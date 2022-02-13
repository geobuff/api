package repo

type ManualTriviaAnswer struct {
	ID                     int    `json:"id"`
	ManualTriviaQuestionID int    `json:"manualTriviaQuestionId"`
	Text                   string `json:"text"`
	IsCorrect              bool   `json:"isCorrect"`
	FlagCode               string `json:"flagCode"`
}

type CreateManualTriviaAnswerDto struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

func CreateManualTriviaAnswer(questionID int, answer CreateManualTriviaAnswerDto) error {
	statement := "INSERT INTO manualtriviaanswers (manualtriviaquestionid, text, iscorrect, flagcode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}

func GetManualTriviaAnswers(questionID int) ([]string, error) {
	rows, err := Connection.Query("SELECT text FROM manualtriviaanswers WHERE manualtriviaquestionid = $1;", questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers = []string{}
	for rows.Next() {
		var answer string
		if err = rows.Scan(&answer); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
