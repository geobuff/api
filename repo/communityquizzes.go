package repo

import (
	"database/sql"
	"time"
)

type CommunityQuiz struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxScore    int       `json:"maxScore"`
	Added       time.Time `json:"added"`
}

type CommunityQuizDto struct {
	ID          int           `json:"id"`
	UserID      int           `json:"userId"`
	Username    string        `json:"username"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	MaxScore    int           `json:"maxScore"`
	Added       time.Time     `json:"added"`
	Plays       sql.NullInt64 `json:"plays"`
}

type GetCommunityQuizDto struct {
	UserID      int                           `json:"userId"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	MaxScore    int                           `json:"maxScore"`
	Questions   []GetCommunityQuizQuestionDto `json:"questions"`
}

type GetCommunityQuizzesFilter struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type CreateCommunityQuizDto struct {
	UserID      int                              `json:"userId"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	MaxScore    int                              `json:"maxScore"`
	Questions   []CreateCommunityQuizQuestionDto `json:"questions"`
}

type UpdateCommunityQuizDto struct {
	UserID      int                              `json:"userId"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	MaxScore    int                              `json:"maxScore"`
	Questions   []UpdateCommunityQuizQuestionDto `json:"questions"`
}

func GetCommunityQuizzes(filter GetCommunityQuizzesFilter) ([]CommunityQuizDto, error) {
	statement := "SELECT q.id, q.userid, u.username, q.name, q.description, q.maxscore, q.added, p.plays FROM communityquizzes q JOIN users u ON u.id = q.userid LEFT JOIN communityquizplays p ON p.communityQuizId = q.id LIMIT $1 OFFSET $2;"
	rows, err := Connection.Query(statement, filter.Limit, filter.Page*filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []CommunityQuizDto{}
	for rows.Next() {
		var quiz CommunityQuizDto
		if err = rows.Scan(&quiz.ID, &quiz.UserID, &quiz.Username, &quiz.Name, &quiz.Description, &quiz.MaxScore, &quiz.Added, &quiz.Plays); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func GetFirstCommunityQuizID(offset int) (int, error) {
	statement := "SELECT id FROM communityquizzes LIMIT 1 OFFSET $1;"
	var id int
	err := Connection.QueryRow(statement, offset).Scan(&id)
	return id, err
}

func GetUserCommunityQuizzes(userID int) ([]CommunityQuizDto, error) {
	statement := "SELECT q.id, q.userid, u.username, q.name, q.description, q.maxscore, q.added, p.plays FROM communityquizzes q JOIN users u ON u.id = q.userid LEFT JOIN communityquizplays p ON p.communityQuizId = q.id WHERE q.userId = $1;"
	rows, err := Connection.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []CommunityQuizDto{}
	for rows.Next() {
		var quiz CommunityQuizDto
		if err = rows.Scan(&quiz.ID, &quiz.UserID, &quiz.Username, &quiz.Name, &quiz.Description, &quiz.MaxScore, &quiz.Added, &quiz.Plays); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func InsertCommunityQuiz(quiz CreateCommunityQuizDto) error {
	statement := "INSERT INTO communityquizzes (userid, name, description, maxscore, added) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	var quizID int
	if err := Connection.QueryRow(statement, quiz.UserID, quiz.Name, quiz.Description, quiz.MaxScore, time.Now()).Scan(&quizID); err != nil {
		return err
	}

	for _, question := range quiz.Questions {
		questionID, err := InsertCommunityQuizQuestion(quizID, question)
		if err != nil {
			return err
		}

		for _, answer := range question.Answers {
			if err := InsertCommunityQuizAnswer(questionID, answer); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateCommunityQuiz(quizID int, quiz UpdateCommunityQuizDto) error {
	statement := "UPDATE communityquizzes SET name = $1, description = $2, maxscore = $3 WHERE id = $4 RETURNING id;"
	var id int
	if err := Connection.QueryRow(statement, quiz.Name, quiz.Description, quiz.MaxScore, quizID).Scan(&id); err != nil {
		return err
	}

	for _, question := range quiz.Questions {
		if question.ID.Valid {
			if err := DeleteCommunityQuizAnswers(int(question.ID.Int64)); err != nil {
				return err
			}

			if err := DeleteCommunityQuizQuestion(int(question.ID.Int64)); err != nil {
				return err
			}
		}

		questionID, err := InsertCommunityQuizQuestion(id, CreateCommunityQuizQuestionDto{
			TypeID:      question.TypeID,
			Question:    question.Question,
			Map:         question.Map,
			Highlighted: question.Highlighted,
			FlagCode:    question.FlagCode,
			ImageUrl:    question.ImageUrl,
		})

		if err != nil {
			return err
		}

		for _, answer := range question.Answers {
			if err := InsertCommunityQuizAnswer(questionID, answer); err != nil {
				return err
			}
		}
	}

	return nil
}

func GetCommunityQuiz(quizID int) (GetCommunityQuizDto, error) {
	statement := "SELECT userid, name, description, maxscore FROM communityquizzes WHERE id = $1;"
	var quiz GetCommunityQuizDto
	if err := Connection.QueryRow(statement, quizID).Scan(&quiz.UserID, &quiz.Name, &quiz.Description, &quiz.MaxScore); err != nil {
		return quiz, err
	}

	questions, err := GetCommunityQuizQuestions(quizID)
	if err != nil {
		return quiz, err
	}
	quiz.Questions = questions

	return quiz, err
}

func DeleteCommunityQuiz(quizID int) error {
	questionIds, err := GetCommunityQuizQuestionIds(quizID)
	if err != nil {
		return err
	}

	for _, questionId := range questionIds {
		if err = DeleteCommunityQuizAnswers(questionId); err != nil {
			return err
		}

		if err = DeleteCommunityQuizQuestion(questionId); err != nil {
			return err
		}
	}

	statement := "DELETE FROM communityquizzes WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, quizID).Scan(&id)
}
