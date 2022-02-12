package repo

type TriviaAnswer struct {
	ID               int    `json:"id"`
	TriviaQuestionID int    `json:"triviaQuestionId"`
	Text             string `json:"text"`
	IsCorrect        bool   `json:"isCorrect"`
	FlagCode         string `json:"flagCode"`
}

type AnswerDto struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

func CreateTriviaAnswer(answer TriviaAnswer) error {
	statement := "INSERT INTO triviaAnswers (triviaQuestionId, text, isCorrect, flagCode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, answer.TriviaQuestionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}
