package repo

import "time"

type CommunityQuiz struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxScore    int       `json:"maxScore"`
	Added       time.Time `json:"added"`
}

type CommunityQuizDto struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxScore    int       `json:"maxScore"`
	Added       time.Time `json:"added"`
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

func GetCommunityQuizzes(filter GetCommunityQuizzesFilter) ([]CommunityQuizDto, error) {
	statement := "SELECT q.id, q.userid, u.username, q.name, q.description, q.maxscore, q.added FROM communityquizzes q JOIN users u ON u.id = q.userid LIMIT $1 OFFSET $2;"
	rows, err := Connection.Query(statement, filter.Limit, filter.Page*filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []CommunityQuizDto{}
	for rows.Next() {
		var quiz CommunityQuizDto
		if err = rows.Scan(&quiz.ID, &quiz.UserID, &quiz.Username, &quiz.Name, &quiz.Description, &quiz.MaxScore, &quiz.Added); err != nil {
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
