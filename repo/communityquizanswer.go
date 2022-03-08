package repo

type CommunityQuizAnswer struct {
	ID                      int    `json:"id"`
	CommunityQuizQuestionID int    `json:"communityQuizQuestionId"`
	Text                    string `json:"text"`
	IsCorrect               bool   `json:"isCorrect"`
	FlagCode                string `json:"flagCode"`
}

type CreateCommunityQuizAnswerDto struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

func InsertCommunityQuizAnswer(questionID int, answer CreateCommunityQuizAnswerDto) error {
	statement := "INSERT INTO communityquizanswers (communityquizquestionid, text, iscorrect, flagcode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}
