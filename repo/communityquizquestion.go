package repo

type CommunityQuizQuestion struct {
	ID              int    `json:"id"`
	CommunityQuizID int    `json:"communityQuizId"`
	TypeID          int    `json:"typeId"`
	Question        string `json:"question"`
	Map             string `json:"map"`
	Highlighted     string `json:"highlighted"`
	FlagCode        string `json:"flagCode"`
	ImageUrl        string `json:"imageUrl"`
}

type CreateCommunityQuizQuestionDto struct {
	TypeID      int                            `json:"typeId"`
	Question    string                         `json:"question"`
	Map         string                         `json:"map"`
	Highlighted string                         `json:"highlighted"`
	FlagCode    string                         `json:"flagCode"`
	ImageUrl    string                         `json:"imageUrl"`
	Answers     []CreateCommunityQuizAnswerDto `json:"answers"`
}

func InsertCommunityQuizQuestion(quizID int, question CreateCommunityQuizQuestionDto) (int, error) {
	statement := "INSERT INTO communityquizquestions (communityquizid, typeid, question, map, highlighted, flagcode, imageurl) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, quizID, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageUrl).Scan(&id)
	return id, err
}
