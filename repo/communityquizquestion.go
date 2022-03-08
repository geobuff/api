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

type UpdateCommunityQuizQuestionDto struct {
	ID          int                            `json:"id"`
	TypeID      int                            `json:"typeId"`
	Question    string                         `json:"question"`
	Map         string                         `json:"map"`
	Highlighted string                         `json:"highlighted"`
	FlagCode    string                         `json:"flagCode"`
	ImageUrl    string                         `json:"imageUrl"`
	Answers     []UpdateCommunityQuizAnswerDto `json:"answers"`
}

func InsertCommunityQuizQuestion(quizID int, question CreateCommunityQuizQuestionDto) (int, error) {
	statement := "INSERT INTO communityquizquestions (communityquizid, typeid, question, map, highlighted, flagcode, imageurl) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, quizID, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageUrl).Scan(&id)
	return id, err
}

func UpdateCommunityQuizQuestion(question UpdateCommunityQuizQuestionDto) error {
	statement := "UPDATE communityquizquestions SET typeid = $1, question = $2, map = $3, highlighted = $4, flagcode = $5, imageurl = $6 WHERE id = $7 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageUrl, question.ID).Scan(&id)
}
