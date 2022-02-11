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
}

type QuestionDto struct {
	ID          int         `json:"id"`
	Type        string      `json:"type"`
	Question    string      `json:"question"`
	Map         string      `json:"map"`
	Highlighted string      `json:"highlighted"`
	FlagCode    string      `json:"flagCode"`
	ImageURL    string      `json:"imageUrl"`
	Answers     []AnswerDto `json:"answers"`
}

func CreateTriviaQuestion(question TriviaQuestion) (int, error) {
	statement := "INSERT INTO triviaQuestions (triviaId, typeId, question, map, highlighted, flagCode, imageUrl) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.TriviaId, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL).Scan(&id)
	return id, err
}
