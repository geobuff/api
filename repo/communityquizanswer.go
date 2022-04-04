package repo

import "database/sql"

type CommunityQuizAnswer struct {
	ID                      int    `json:"id"`
	CommunityQuizQuestionID int    `json:"communityQuizQuestionId"`
	Text                    string `json:"text"`
	IsCorrect               bool   `json:"isCorrect"`
	FlagCode                string `json:"flagCode"`
}

type GetCommunityQuizAnswerDto struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

type CreateCommunityQuizAnswerDto struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

type UpdateCommunityQuizAnswerDto struct {
	ID        sql.NullInt64 `json:"id"`
	Text      string        `json:"text"`
	IsCorrect bool          `json:"isCorrect"`
	FlagCode  string        `json:"flagCode"`
}

func GetCommunityQuizAnswers(questionID int) ([]GetCommunityQuizAnswerDto, error) {
	statement := "SELECT id, text, iscorrect, flagcode FROM communityquizanswers WHERE communityquizquestionid = $1;"
	rows, err := Connection.Query(statement, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers = []GetCommunityQuizAnswerDto{}
	for rows.Next() {
		var answer GetCommunityQuizAnswerDto
		if err = rows.Scan(&answer.ID, &answer.Text, &answer.IsCorrect, &answer.FlagCode); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, rows.Err()
}

func InsertCommunityQuizAnswer(questionID int, answer CreateCommunityQuizAnswerDto) error {
	statement := "INSERT INTO communityquizanswers (communityquizquestionid, text, iscorrect, flagcode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}

func UpdateCommunityQuizAnswer(answerID int, answer UpdateCommunityQuizAnswerDto) error {
	statement := "UPDATE communityquizanswers SET text = $1, iscorrect = $2, flagcode = $3 WHERE id = $4 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, answer.Text, answer.IsCorrect, answer.FlagCode, answerID).Scan(&id)
}

func DeleteCommunityQuizAnswers(questionID int) error {
	statement := "DELETE FROM communityquizanswers WHERE communityquizquestionid = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, questionID).Scan(&id)
}
