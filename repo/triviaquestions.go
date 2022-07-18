package repo

type TriviaQuestion struct {
	ID                 int    `json:"id"`
	TriviaId           int    `json:"triviaId"`
	TypeID             int    `json:"typeId"`
	Question           string `json:"question"`
	Map                string `json:"map"`
	Highlighted        string `json:"highlighted"`
	FlagCode           string `json:"flagCode"`
	ImageURL           string `json:"imageUrl"`
	ImageAttributeName string `json:"imageAttributeName"`
	ImageAttributeURL  string `json:"ImageAttributeUrl"`
	ImageWidth         int    `json:"imageWidth"`
	ImageHeight        int    `json:"imageHeight"`
	ImageAlt           string `json:"imageAlt"`
	Explainer          string `json:"explainer"`
}

type QuestionDto struct {
	ID                 int         `json:"id"`
	Type               string      `json:"type"`
	Question           string      `json:"question"`
	Map                string      `json:"map"`
	Highlighted        string      `json:"highlighted"`
	FlagCode           string      `json:"flagCode"`
	ImageURL           string      `json:"imageUrl"`
	ImageAttributeName string      `json:"imageAttributeName"`
	ImageAttributeURL  string      `json:"imageAttributeUrl"`
	ImageWidth         int         `json:"imageWidth"`
	ImageHeight        int         `json:"imageHeight"`
	ImageAlt           string      `json:"imageAlt"`
	Explainer          string      `json:"explainer"`
	Answers            []AnswerDto `json:"answers"`
}

func CreateTriviaQuestion(question TriviaQuestion) (int, error) {
	statement := "INSERT INTO triviaQuestions (triviaId, typeId, question, map, highlighted, flagCode, imageUrl, imageAttributeName, imageAttributeUrl, imageWidth, imageHeight, imageAlt, explainer) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.TriviaId, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL, question.ImageAttributeName, question.ImageAttributeURL, question.ImageWidth, question.ImageHeight, question.ImageAlt, question.Explainer).Scan(&id)
	return id, err
}

func DeleteTriviaQuestion(questionId int) error {
	statement := "DELETE FROM triviaQuestions WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionId).Scan(&id)
}
