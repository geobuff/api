package repo

type TriviaQuestion struct {
	ID          int    `json:"id"`
	TriviaId    int    `json:"triviaId"`
	TypeID      int    `json:"typeId"`
	Question    string `json:"question"`
	Map         string `json:"map"`
	Highlighted string `json:"highlighted"`
	FlagCode    string `json:"flagCode"`
	ImageURL    string `json:"imageUrl"`
	Explainer   string `json:"explainer"`
}

type QuestionDto struct {
	ID          int         `json:"id"`
	Type        string      `json:"type"`
	Question    string      `json:"question"`
	Map         string      `json:"map"`
	Highlighted string      `json:"highlighted"`
	FlagCode    string      `json:"flagCode"`
	ImageURL    string      `json:"imageUrl"`
	Explainer   string      `json:"explainer"`
	Answers     []AnswerDto `json:"answers"`
}

func CreateTriviaQuestion(question TriviaQuestion) (int, error) {
	statement := "INSERT INTO triviaQuestions (triviaId, typeId, question, map, highlighted, flagCode, imageUrl, explainer) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.TriviaId, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL, question.Explainer).Scan(&id)
	return id, err
}

func DeleteTriviaQuestion(questionId int) error {
	statement := "DELETE FROM triviaQuestions WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionId).Scan(&id)
}
